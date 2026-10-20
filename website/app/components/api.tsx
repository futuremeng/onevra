import { useI18n } from "@fumadocs/base-ui/contexts/i18n"
import { ExternalLink } from "lucide-react"
import type { ReactNode } from "react"

const text = {
  en: { methods: "Methods", reference: "pkg.go.dev" },
  zh: { methods: "方法", reference: "pkg.go.dev" },
}

/**
 * One exported function or type, which <GoAPI> expands into at compile time
 * (see scripts/remark-onevra.ts): its declaration, doc comment, fields,
 * values and methods.
 */
export function APIEntry({
  href,
  children,
}: {
  kind: "func" | "type"
  href: string
  children?: ReactNode
}) {
  const t = useI18n().locale === "zh" ? text.zh : text.en
  return (
    <div className="api-entry relative my-4 [&>figure:first-child]:mt-0">
      {children}
      <a
        href={href}
        target="_blank"
        rel="noreferrer"
        className="not-prose inline-flex items-center gap-1 text-xs text-fd-muted-foreground no-underline hover:text-fd-foreground"
      >
        {t.reference}
        <ExternalLink className="size-3" aria-hidden />
      </a>
    </div>
  )
}

/** A type's methods, folded away until opened. */
export function APIMethods({ children }: { children?: ReactNode }) {
  const t = useI18n().locale === "zh" ? text.zh : text.en
  return (
    <details className="group my-4 rounded-xl border bg-fd-card px-4 [&_figure]:my-3">
      <summary className="cursor-pointer py-3 text-sm font-medium text-fd-foreground select-none">
        {t.methods}
      </summary>
      <div className="pb-1">{children}</div>
    </details>
  )
}
