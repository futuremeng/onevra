// Fills two MDX elements of the documentation from the Go source, at compile
// time, so that pages carry only what they show and code is highlighted with
// the rest of the page:
//
//   <Example name="switch" />            an example of source/examples/website,
//                                        given its source as a code block
//   <GoAPI name="components.Switch" />   a heading and the declaration, doc
//                                        comment, fields, values and methods
//
// The data comes from .generated/docs-data.json (see examples.ts). An unknown
// name fails the build.
import { readFileSync, statSync } from "node:fs"
import { valueToEstree } from "estree-util-value-to-estree"
import type { Code, Heading, Root, RootContent } from "mdast"
import { fromMarkdown } from "mdast-util-from-markdown"
import { gfmFromMarkdown } from "mdast-util-gfm"
import type { MdxJsxAttribute, MdxJsxFlowElement } from "mdast-util-mdx-jsx"
import { gfm } from "micromark-extension-gfm"
import { visit } from "unist-util-visit"

import { dataPath } from "./examples.ts"

interface ExampleData {
  name: string
  file: string
  code: string
}

interface APIData {
  kind: "func" | "type"
  package: string
  decl: string
  doc?: string
  fields?: { name: string; type: string; doc?: string }[]
  values?: { decl: string; doc?: string }[]
  methods?: { name: string; decl: string; doc?: string }[]
}

interface DocsData {
  examples: ExampleData[]
  api: Record<string, APIData>
}

let cache: { mtime: number; data: DocsData } | undefined

function load(): DocsData {
  const mtime = statSync(dataPath).mtimeMs
  if (cache?.mtime !== mtime) {
    cache = { mtime, data: JSON.parse(readFileSync(dataPath, "utf8")) }
  }
  return cache.data
}

const repository = "https://github.com/oneo-me/onevra/blob/main/"

export function remarkOnevra() {
  return (tree: Root) => {
    const data = load()
    visit(tree, "mdxJsxFlowElement", (node, index, parent) => {
      if (node.name === "Example") {
        const name = attribute(node, "name")
        const example = data.examples.find((e) => e.name === name)
        if (!example) {
          throw new Error(
            `<Example name="${name}">: no example of that name in source/examples/website/examples`
          )
        }
        node.attributes.push(
          stringAttribute("source", repository + example.file)
        )
        node.children = [
          code(example.code, `title="${name.replaceAll("-", "_")}.go"`),
        ]
        return "skip"
      }
      if (node.name === "GoAPI" && parent && index !== undefined) {
        const name = attribute(node, "name")
        const api = data.api[name]
        if (!api) {
          throw new Error(
            `<GoAPI name="${name}">: no exported API of that name`
          )
        }
        parent.children.splice(index, 1, ...apiNodes(name, api))
        return "skip"
      }
    })
  }
}

function attribute(node: MdxJsxFlowElement, name: string): string {
  for (const a of node.attributes) {
    if (
      a.type === "mdxJsxAttribute" &&
      a.name === name &&
      typeof a.value === "string"
    ) {
      return a.value
    }
  }
  throw new Error(`<${node.name}> needs a ${name} attribute`)
}

function stringAttribute(name: string, value: string): MdxJsxAttribute {
  return { type: "mdxJsxAttribute", name, value }
}

function expressionAttribute(name: string, value: unknown): MdxJsxAttribute {
  return {
    type: "mdxJsxAttribute",
    name,
    value: {
      type: "mdxJsxAttributeValueExpression",
      value: JSON.stringify(value),
      data: {
        estree: {
          type: "Program",
          sourceType: "module",
          body: [
            { type: "ExpressionStatement", expression: valueToEstree(value) },
          ],
        },
      },
    },
  }
}

function code(value: string, meta?: string): Code {
  return { type: "code", lang: "go", meta, value: value.replace(/\n$/, "") }
}

function markdown(text: string | undefined): RootContent[] {
  if (!text) return []
  const children = fromMarkdown(text, {
    extensions: [gfm()],
    mdastExtensions: [gfmFromMarkdown()],
  }).children
  // Positions in the doc comment are not positions in the page, which the
  // search index slices its text from; the entry carries its own instead.
  visit({ type: "root", children } as Root, (node) => {
    delete node.position
  })
  return children
}

function element(
  name: string,
  attributes: MdxJsxAttribute[],
  children: RootContent[]
): MdxJsxFlowElement {
  return {
    type: "mdxJsxFlowElement",
    name,
    attributes,
    children: children as MdxJsxFlowElement["children"],
  }
}

function apiNodes(key: string, api: APIData): RootContent[] {
  const name = key.slice(key.indexOf(".") + 1)
  const heading: Heading = {
    type: "heading",
    depth: 3,
    children: [{ type: "inlineCode", value: name }],
  }
  const children: RootContent[] = []
  const struct = api.fields !== undefined && api.fields.length > 0
  if (!struct) children.push(code(api.decl))
  children.push(...markdown(api.doc))
  if (struct) {
    const type: Record<
      string,
      { type: string; description?: string; required: boolean }
    > = {}
    for (const f of api.fields!) {
      // Go has no optional fields to mark; every field is set or zero.
      type[f.name] = {
        type: f.type,
        description: f.doc || undefined,
        required: true,
      }
    }
    children.push(element("TypeTable", [expressionAttribute("type", type)], []))
  }
  for (const v of api.values ?? []) {
    children.push(...markdown(v.doc), code(v.decl))
  }
  if (api.methods?.length) {
    children.push(
      element(
        "APIMethods",
        [],
        api.methods.flatMap((m) => [code(m.decl), ...markdown(m.doc)])
      )
    )
  }
  const entry = element(
    "APIEntry",
    [
      stringAttribute("kind", api.kind),
      stringAttribute("href", `https://pkg.go.dev/${api.package}#${name}`),
    ],
    children
  )
  // Search finds the entry by its name and doc comment.
  entry.data = {
    structuredData: {
      contents: [
        {
          heading: name.toLowerCase(),
          content: `${name}: ${(api.doc ?? "").replace(/\s+/g, " ")}`,
        },
      ],
    },
  } as MdxJsxFlowElement["data"]
  return [heading, entry]
}
