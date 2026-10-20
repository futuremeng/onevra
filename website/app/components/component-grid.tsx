import { useTreeContext } from "@fumadocs/base-ui/contexts/tree"
import type * as PageTree from "fumadocs-core/page-tree"
import type { ReactNode } from "react"
import { Link } from "react-router"

interface Group {
  name: ReactNode
  pages: PageTree.Item[]
}

/** Groups the pages of the current sidebar root under its separators. */
function groups(nodes: PageTree.Node[]): Group[] {
  const out: Group[] = []
  for (const node of nodes) {
    if (node.type === "separator") {
      out.push({ name: node.name, pages: [] })
    } else if (node.type === "page") {
      if (out.length === 0) continue // the overview itself
      out[out.length - 1].pages.push(node)
    } else if (node.type === "folder") {
      out.push({ name: node.name, pages: [] })
      for (const child of node.children) {
        if (child.type === "page") out[out.length - 1].pages.push(child)
      }
    }
  }
  return out.filter((g) => g.pages.length > 0)
}

/**
 * Every component page of the sidebar, by group, so that the overview
 * follows meta.json as pages come and go.
 */
export function ComponentGrid() {
  const { root } = useTreeContext()
  return (
    <div className="not-prose flex flex-col gap-10">
      {groups(root.children).map((group, i) => (
        <section key={i}>
          <h2 className="mb-3 text-sm font-medium text-fd-muted-foreground">
            {group.name}
          </h2>
          <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
            {group.pages.map((page) => (
              <Link
                key={page.url}
                to={page.url}
                prefetch="intent"
                className="group flex flex-col gap-1 rounded-xl border bg-fd-card p-4 transition-colors hover:border-fd-primary/40 hover:bg-fd-accent/60"
              >
                <span className="flex items-center gap-2 font-medium text-fd-card-foreground">
                  {page.icon ? (
                    <span className="text-fd-primary [&_svg]:size-4">
                      {page.icon}
                    </span>
                  ) : null}
                  {page.name}
                </span>
                {page.description ? (
                  <span className="text-sm text-fd-muted-foreground">
                    {page.description}
                  </span>
                ) : null}
              </Link>
            ))}
          </div>
        </section>
      ))}
    </div>
  )
}
