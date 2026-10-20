import { readFile } from "node:fs/promises"
import { resolve } from "node:path"

import { codeToHtml } from "shiki"
import type { Locale } from "~/lib/i18n"
import { source } from "~/lib/source"

// Loaders run while the site builds, bundled away from this file, so paths
// start from the website directory the build runs in.
const dataPath = resolve(".generated/docs-data.json")
const changelogPath = resolve("../CHANGELOG.md")

/** The latest released version, from the changelog's first versioned heading. */
async function version(): Promise<string> {
  const changelog = await readFile(changelogPath, "utf8")
  return /^## \[(\d[^\]]*)\]/m.exec(changelog)?.[1] ?? "latest"
}

/** An example's source as highlighted HTML, in both themes. */
async function exampleHTML(name: string): Promise<string> {
  const data = JSON.parse(await readFile(dataPath, "utf8")) as {
    examples: { name: string; code: string }[]
  }
  const example = data.examples.find((e) => e.name === name)
  if (!example) throw new Error(`no example named ${name}`)
  return codeToHtml(example.code.trimEnd(), {
    lang: "go",
    themes: { light: "github-light", dark: "github-dark" },
    defaultColor: false,
  })
}

export interface ComponentLink {
  title: string
  url: string
}

/** The component pages, in sidebar order, without the overview. */
function components(locale: Locale): ComponentLink[] {
  return source
    .getPages(locale)
    .filter((p) => /\/docs\/components\/.+/.test(p.url))
    .map((p) => ({ title: p.data.title, url: p.url }))
}

export async function homeData(locale: Locale) {
  return {
    version: await version(),
    todo: await exampleHTML("todo"),
    components: components(locale),
  }
}
