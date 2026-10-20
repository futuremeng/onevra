// Builds the documentation's live examples from source/examples/website:
// writes the examples' source and the Go API to .generated/docs-data.json,
// which remark-onevra.ts reads while compiling MDX, and compiles the examples
// to WebAssembly, served in development and emitted into the build. Go must be
// installed; the build fails without it.
import { execFile } from "node:child_process"
import { createHash } from "node:crypto"
import { mkdir, mkdtemp, readFile, rm, writeFile } from "node:fs/promises"
import { tmpdir } from "node:os"
import { dirname, join, resolve } from "node:path"
import { fileURLToPath } from "node:url"
import { promisify } from "node:util"
import type { Plugin, ViteDevServer } from "vite"

const run = promisify(execFile)

const websiteDir = fileURLToPath(new URL("../", import.meta.url))
const sourceDir = resolve(websiteDir, "../source")
export const dataPath = resolve(websiteDir, ".generated/docs-data.json")
const goEnv = { ...process.env, GOWORK: "off" }

/** Writes the examples' source and the API for the MDX compiler. */
export async function generateData(): Promise<void> {
  try {
    const { stdout } = await run(
      "go",
      ["run", "./examples/website", "-manifest"],
      { cwd: sourceDir, env: goEnv, maxBuffer: 1 << 26 }
    )
    await mkdir(dirname(dataPath), { recursive: true })
    await writeFile(dataPath, stdout)
  } catch (error) {
    const e = error as { stderr?: string; message: string }
    throw new Error(
      `onevra-examples: go run ./examples/website -manifest failed: ${e.stderr || e.message}`
    )
  }
}

interface Build {
  wasm: Buffer
  loader: string
  name: string
}

async function compile(): Promise<Build> {
  const dir = await mkdtemp(join(tmpdir(), "onevra-examples-"))
  const out = join(dir, "app.wasm")
  try {
    await run(
      "go",
      [
        "build",
        "-tags=mygo_noinspector",
        "-trimpath",
        "-ldflags=-s -w",
        "-o",
        out,
        "./examples/website",
      ],
      {
        cwd: sourceDir,
        env: { ...goEnv, GOOS: "js", GOARCH: "wasm", CGO_ENABLED: "0" },
        maxBuffer: 1 << 24,
      }
    )
    const { stdout: goroot } = await run("go", ["env", "GOROOT"], {
      cwd: sourceDir,
    })
    const [wasm, loader] = await Promise.all([
      readFile(out),
      readFile(join(goroot.trim(), "lib", "wasm", "wasm_exec.js"), "utf8"),
    ])
    // Named by its content, so that browsers keep it while it lasts.
    const hash = createHash("sha256").update(wasm).digest("hex").slice(0, 12)
    return { wasm, loader, name: `app.${hash}.wasm` }
  } catch (error) {
    const e = error as { stderr?: string; message: string }
    throw new Error(
      `onevra-examples: go build of ./examples/website failed: ${e.stderr || e.message}`
    )
  } finally {
    await rm(dir, { recursive: true, force: true })
  }
}

// Shared by the client and server builds, which run in one process.
let building: Promise<Build> | undefined
let generating: Promise<void> | undefined

function build(): Promise<Build> {
  building ??= compile().catch((error) => {
    building = undefined
    throw error
  })
  return building
}

function generate(): Promise<void> {
  generating ??= generateData().catch((error) => {
    generating = undefined
    throw error
  })
  return generating
}

const virtualId = "virtual:onevra-examples"
const resolvedId = "\0" + virtualId

export function onevraExamples(): Plugin {
  let base = "/"
  let dev = false
  return {
    name: "onevra-examples",
    configResolved(config) {
      base = config.base
      dev = config.command === "serve"
    },
    async buildStart() {
      await generate()
    },
    resolveId(id) {
      return id === virtualId ? resolvedId : undefined
    },
    async load(id) {
      if (id !== resolvedId) return
      // Development serves a fixed name with no-store; the build names the
      // file by its content.
      const name = dev ? "app.wasm" : (await build()).name
      return [
        `export const wasmUrl = ${JSON.stringify(`${base}examples/${name}`)}`,
        `export const loaderUrl = ${JSON.stringify(`${base}examples/wasm_exec.js`)}`,
      ].join("\n")
    },
    async generateBundle() {
      if (this.environment.name !== "client") return
      const { wasm, loader, name } = await build()
      this.emitFile({
        type: "asset",
        fileName: `examples/${name}`,
        source: wasm,
      })
      this.emitFile({
        type: "asset",
        fileName: "examples/wasm_exec.js",
        source: loader,
      })
    },
    configureServer(server) {
      watchSource(server)
      server.middlewares.use(async (request, response, next) => {
        const { pathname } = new URL(request.url ?? "/", "http://localhost")
        const file = pathname.slice(`${base}examples/`.length)
        if (!pathname.startsWith(`${base}examples/`)) return next()
        try {
          if (file === "app.wasm") {
            response.setHeader("Content-Type", "application/wasm")
            response.setHeader("Cache-Control", "no-store")
            response.end((await build()).wasm)
          } else if (file === "wasm_exec.js") {
            response.setHeader("Content-Type", "text/javascript; charset=utf-8")
            response.setHeader("Cache-Control", "no-store")
            response.end((await build()).loader)
          } else {
            next()
          }
        } catch (error) {
          next(error)
        }
      })
    },
  }
}

/**
 * Rebuilds after Go source changes: the next load of an example gets the new
 * WebAssembly, and pages compile again from the new data.
 */
function watchSource(server: ViteDevServer) {
  server.watcher.add(sourceDir)
  let timer: ReturnType<typeof setTimeout> | undefined
  server.watcher.on("all", (event, file) => {
    if (!file.startsWith(sourceDir) || !file.endsWith(".go")) return
    if (event !== "change" && event !== "add" && event !== "unlink") return
    clearTimeout(timer)
    timer = setTimeout(async () => {
      building = undefined
      generating = undefined
      try {
        await generate()
      } catch (error) {
        server.config.logger.error(String(error))
        return
      }
      for (const environment of Object.values(server.environments)) {
        for (const module of environment.moduleGraph.idToModuleMap.values()) {
          if (module.id && /\.mdx($|\?)/.test(module.id)) {
            environment.moduleGraph.invalidateModule(module)
          }
        }
      }
      server.ws.send({ type: "full-reload" })
    }, 200)
  })
}
