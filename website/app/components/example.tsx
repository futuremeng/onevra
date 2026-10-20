import { ChevronDown, ExternalLink, LoaderCircle } from "lucide-react"
import { useI18n } from "@fumadocs/base-ui/contexts/i18n"
import { useEffect, useRef, useState, type ReactNode } from "react"

import { mountExample, setExampleVisible, unmountExample } from "~/lib/examples"
import { cn } from "~/lib/utils"

type Status = "idle" | "loading" | "live" | "error"

const text = {
  en: {
    loading: "Loading the live example…",
    failed: "The live example did not load.",
    retry: "Try again",
    source: "Source",
    expand: "Show all code",
    collapse: "Show less",
  },
  zh: {
    loading: "正在加载实时示例…",
    failed: "实时示例加载失败。",
    retry: "重试",
    source: "源码",
    expand: "展开全部代码",
    collapse: "收起",
  },
}

function useLocaleText() {
  return useI18n().locale === "zh" ? text.zh : text.en
}

/** Code taller than this collapses behind a button. */
const collapsedHeight = 360

/**
 * A live example of source/examples/website/examples, running in the page,
 * above its source. The MDX compiler gives it the source as its children and
 * its file's address as source (see scripts/remark-onevra.ts).
 */
export function Example({
  name,
  height = 220,
  scroll = false,
  source,
  children,
}: {
  name: string
  height?: number
  scroll?: boolean
  source?: string
  children?: ReactNode
}) {
  const t = useLocaleText()
  return (
    <figure className="not-prose my-6 overflow-hidden rounded-xl border bg-fd-card shadow-xs">
      <LiveCanvas name={name} scroll={scroll} style={{ height }} />
      <ExampleCode source={source} labels={t}>
        {children}
      </ExampleCode>
    </figure>
  )
}

/**
 * The window of a live example. It opens when it nears the viewport and may
 * close far from it, as browsers keep only a few WebGL contexts. Until it is
 * clicked, the wheel scrolls the page rather than the example; set scroll to
 * give it the wheel at once.
 */
export function LiveCanvas({
  name,
  scroll = false,
  className,
  style,
}: {
  name: string
  scroll?: boolean
  className?: string
  style?: React.CSSProperties
}) {
  const t = useLocaleText()
  const frame = useRef<HTMLDivElement>(null)
  const target = useRef<HTMLDivElement>(null)
  const [status, setStatus] = useState<Status>("idle")
  const [attempt, setAttempt] = useState(0)
  const active = useRef(scroll)

  useEffect(() => {
    const element = target.current
    if (!element) return
    let disposed = false
    let current: Awaited<ReturnType<typeof mountExample>> | undefined
    let mounting = false
    let visible = false

    const open = () => {
      if (mounting || current || disposed) return
      mounting = true
      setStatus("loading")
      mountExample(element, name).then(
        (w) => {
          mounting = false
          if (disposed) return unmountExample(w)
          current = w
          setExampleVisible(w, visible)
          setStatus("live")
        },
        (error) => {
          mounting = false
          console.error(error)
          if (!disposed) setStatus("error")
        }
      )
    }
    const evicted = () => {
      current = undefined
      setStatus("idle")
    }
    element.addEventListener("onevra-example-evicted", evicted)

    const observer = new IntersectionObserver(
      ([entry]) => {
        visible = entry.isIntersecting
        if (current) setExampleVisible(current, visible)
        else if (visible) open()
      },
      { rootMargin: "400px 0px" }
    )
    observer.observe(element)

    return () => {
      disposed = true
      observer.disconnect()
      element.removeEventListener("onevra-example-evicted", evicted)
      if (current) unmountExample(current)
    }
  }, [name, attempt])

  // MyGo takes every wheel event over its canvas; until the example is
  // clicked, stop them short of it so that the page scrolls past.
  useEffect(() => {
    const element = frame.current
    if (!element) return
    const wheel = (e: WheelEvent) => {
      if (!active.current) e.stopPropagation()
    }
    element.addEventListener("wheel", wheel, { capture: true })
    return () => element.removeEventListener("wheel", wheel, { capture: true })
  }, [])

  return (
    <div
      ref={frame}
      className={cn("relative bg-fd-background", className)}
      style={style}
      onPointerDown={() => (active.current = true)}
      onPointerLeave={() => (active.current = scroll)}
    >
      <div ref={target} className="absolute inset-0" />
      {status !== "live" ? (
        <div className="absolute inset-0 flex flex-col items-center justify-center gap-2 text-sm text-fd-muted-foreground">
          {status === "error" ? (
            <>
              <p>{t.failed}</p>
              <button
                type="button"
                onClick={() => setAttempt((n) => n + 1)}
                className="rounded-md border px-3 py-1 text-fd-foreground hover:bg-fd-accent"
              >
                {t.retry}
              </button>
            </>
          ) : (
            <>
              <LoaderCircle className="size-5 animate-spin" aria-hidden />
              <p>{t.loading}</p>
            </>
          )}
        </div>
      ) : null}
    </div>
  )
}

function ExampleCode({
  source,
  labels,
  children,
}: {
  source?: string
  labels: (typeof text)["en"]
  children?: ReactNode
}) {
  const body = useRef<HTMLDivElement>(null)
  const [tall, setTall] = useState(false)
  const [expanded, setExpanded] = useState(false)

  useEffect(() => {
    const element = body.current
    if (!element) return
    const measure = () => setTall(element.scrollHeight > collapsedHeight + 48)
    measure()
    const observer = new ResizeObserver(measure)
    observer.observe(element)
    return () => observer.disconnect()
  }, [])

  return (
    <div className="relative border-t">
      <div
        ref={body}
        className={cn(
          "overflow-hidden [&_figure]:my-0 [&_figure]:rounded-none [&_figure]:border-0 [&_figure]:shadow-none",
          tall && !expanded && "max-h-(--collapsed)"
        )}
        style={{ "--collapsed": `${collapsedHeight}px` } as React.CSSProperties}
      >
        {children}
      </div>
      {source ? (
        <a
          href={source}
          target="_blank"
          rel="noreferrer"
          className="absolute top-1.5 right-11 z-10 inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs text-fd-muted-foreground hover:bg-fd-accent hover:text-fd-accent-foreground"
        >
          {labels.source}
          <ExternalLink className="size-3" aria-hidden />
        </a>
      ) : null}
      {tall ? (
        <div
          className={cn(
            "flex justify-center p-2",
            !expanded &&
              "absolute inset-x-0 bottom-0 bg-linear-to-t from-fd-card via-fd-card/90 to-transparent pt-12"
          )}
        >
          <button
            type="button"
            onClick={() => setExpanded((v) => !v)}
            className="inline-flex items-center gap-1 rounded-full border bg-fd-background px-3 py-1 text-xs font-medium text-fd-foreground shadow-xs hover:bg-fd-accent"
          >
            {expanded ? labels.collapse : labels.expand}
            <ChevronDown
              className={cn(
                "size-3 transition-transform",
                expanded && "rotate-180"
              )}
              aria-hidden
            />
          </button>
        </div>
      ) : null}
    </div>
  )
}
