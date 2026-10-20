// Checks that the documentation covers the Go source: every exported
// function of the documented packages has a <GoAPI> section, every example
// is shown by a page, and every English page has a Chinese translation.
//
//   bun run check
import { readFile, glob } from "node:fs/promises"

import { dataPath, generateData } from "./examples.ts"

await generateData()
const data = JSON.parse(await readFile(dataPath, "utf8")) as {
  examples: { name: string }[]
  api: Record<string, { kind: string }>
}

const apis = new Set<string>()
const examples = new Set<string>()
const pages = new Set<string>()
for await (const file of glob("**/*.mdx", { cwd: "content/docs" })) {
  pages.add(file)
  const text = await readFile(`content/docs/${file}`, "utf8")
  for (const m of text.matchAll(/<GoAPI name="([^"]+)"/g)) apis.add(m[1])
  for (const m of text.matchAll(/<Example name="([^"]+)"/g)) examples.add(m[1])
}
// The home page shows these itself.
for (const name of ["showcase", "todo"]) examples.add(name)

const problems: string[] = []
for (const [name, entry] of Object.entries(data.api)) {
  if (entry.kind === "func" && !apis.has(name))
    problems.push(`no <GoAPI name="${name}" /> on any page`)
}
for (const { name } of data.examples) {
  if (!examples.has(name)) problems.push(`example ${name} is on no page`)
}
for (const page of pages) {
  if (
    !page.endsWith(".zh.mdx") &&
    !pages.has(page.replace(/\.mdx$/, ".zh.mdx"))
  ) {
    problems.push(`${page} has no Chinese translation`)
  }
}
if (problems.length > 0) {
  console.error(problems.map((p) => `check-docs: ${p}`).join("\n"))
  process.exit(1)
}
console.log("check-docs: the documentation covers the API and every example")
