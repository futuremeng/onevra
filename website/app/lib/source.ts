import { loader } from "fumadocs-core/source"
import { defineDocs } from "fumadocs-mdx/macro"
import { Blocks, BookOpen, type LucideIcon } from "lucide-react"
import { createElement } from "react"

import { i18n } from "~/lib/i18n"

export const docs = defineDocs({
  dir: "content/docs",
  docs: { async: true },
})

// The icons meta.json files may name; listed one by one, so the bundle
// carries only these.
const icons: Record<string, LucideIcon> = { Blocks, BookOpen }

export const source = loader({
  source: docs.toFumadocsSource(),
  baseUrl: "/docs",
  i18n,
  icon(name) {
    const icon = name ? icons[name] : undefined
    return icon ? createElement(icon) : undefined
  },
})
