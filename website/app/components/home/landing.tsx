import { HomeLayout } from "@fumadocs/base-ui/layouts/home"
import {
  ArrowRight,
  BarChart3,
  Check,
  FlaskConical,
  MessagesSquare,
  PanelsTopLeft,
  Shapes,
  SlidersHorizontal,
} from "lucide-react"
import { Link } from "react-router"

import galleryUrl from "../../../../docs/images/gallery.png"
import { LiveCanvas } from "~/components/example"
import { GithubIcon } from "~/components/github-icon"
import { InstallCommand } from "~/components/home/install-command"
import { homeText } from "~/components/home/text"
import { WindowFrame } from "~/components/home/window-frame"
import type { ComponentLink } from "~/lib/home.server"
import { localized, type Locale } from "~/lib/i18n"
import { baseOptions } from "~/lib/layout.shared"
import { site } from "~/lib/site"
import { cn } from "~/lib/utils"

const featureIcons = [
  SlidersHorizontal,
  PanelsTopLeft,
  BarChart3,
  MessagesSquare,
  Shapes,
  FlaskConical,
]

const primaryLink =
  "inline-flex h-11 items-center gap-2 rounded-full bg-fd-primary px-6 text-sm font-medium text-fd-primary-foreground shadow-sm transition-colors hover:bg-fd-primary/90"
const secondaryLink =
  "inline-flex h-11 items-center gap-2 rounded-full border bg-fd-background px-6 text-sm font-medium transition-colors hover:bg-fd-accent"

function Section({
  title,
  lead,
  className,
  children,
}: {
  title: string
  lead: string
  className?: string
  children: React.ReactNode
}) {
  return (
    <section className={cn("border-t", className)}>
      <div className="mx-auto w-full max-w-6xl px-6 py-24">
        <div className="mx-auto mb-12 max-w-2xl text-center">
          <h2 className="text-3xl font-semibold tracking-tight text-balance sm:text-4xl">
            {title}
          </h2>
          <p className="mt-4 text-lg text-pretty text-fd-muted-foreground">
            {lead}
          </p>
        </div>
        {children}
      </div>
    </section>
  )
}

export function Landing({
  locale,
  version,
  todo,
  components,
}: {
  locale: Locale
  version: string
  todo: string
  components: ComponentLink[]
}) {
  const t = homeText[locale]
  const install = `go get ${site.module}@${version === "latest" ? "latest" : `v${version}`}`
  const quickStart = localized(locale, "/docs/quick-start")
  const allComponents = localized(locale, "/docs/components")

  return (
    <HomeLayout {...baseOptions(locale)} className="overflow-x-clip">
      {/* Hero */}
      <section className="relative isolate">
        <div
          aria-hidden
          className="absolute inset-0 -z-10 bg-[radial-gradient(ellipse_80%_60%_at_50%_0%,color-mix(in_oklab,var(--color-fd-primary)_14%,transparent),transparent_70%)]"
        />
        <div
          aria-hidden
          className="absolute inset-0 -z-10 bg-[radial-gradient(var(--color-fd-border)_1px,transparent_1px)] [mask-image:radial-gradient(ellipse_70%_55%_at_50%_20%,black,transparent)] bg-size-[22px_22px] opacity-70"
        />
        <div className="mx-auto flex max-w-6xl flex-col items-center px-6 pt-20 pb-10 text-center sm:pt-28">
          <a
            href={site.changelog}
            target="_blank"
            rel="noreferrer"
            className="mb-8 inline-flex items-center gap-2 rounded-full border bg-fd-card/80 px-3 py-1 text-xs font-medium text-fd-muted-foreground backdrop-blur transition-colors hover:text-fd-foreground"
          >
            <span className="rounded-full bg-fd-primary/12 px-2 py-0.5 text-fd-primary">
              v{version}
            </span>
            {t.badge}
            <ArrowRight className="size-3" />
          </a>
          <h1 className="text-5xl font-semibold tracking-tight text-balance sm:text-7xl">
            {t.headline[0]}
            <br />
            <span className="bg-linear-to-r from-fd-primary via-[color-mix(in_oklab,var(--color-fd-primary)_70%,#f5b14c)] to-fd-primary bg-clip-text text-transparent">
              {t.headline[1]}
            </span>
          </h1>
          <p className="mt-6 max-w-2xl text-lg text-pretty text-fd-muted-foreground sm:text-xl">
            {t.lead}
          </p>
          <div className="mt-10 flex flex-wrap justify-center gap-3">
            <Link to={quickStart} className={primaryLink}>
              {t.start}
              <ArrowRight className="size-4" />
            </Link>
            <Link to={allComponents} className={secondaryLink}>
              {t.browse}
            </Link>
          </div>
          <div className="mt-6 max-w-full">
            <InstallCommand command={install} copy={t.copy} copied={t.copied} />
          </div>
        </div>
        <div className="mx-auto max-w-6xl px-4 pb-6 sm:px-6">
          <WindowFrame title={t.live}>
            <LiveCanvas name="showcase" className="h-[500px]" />
          </WindowFrame>
          <p className="mx-auto mt-5 max-w-xl text-center text-sm text-fd-muted-foreground">
            {t.liveNote}
          </p>
        </div>
      </section>

      {/* Code and its result */}
      <Section title={t.code.title} lead={t.code.lead} className="border-t-0">
        <div className="grid items-stretch gap-6 lg:grid-cols-2">
          <div className="flex min-w-0 flex-col overflow-hidden rounded-2xl border bg-fd-card">
            <div className="flex h-10 items-center border-b px-4 font-mono text-xs text-fd-muted-foreground">
              todo.go
            </div>
            <div
              className="max-h-[560px] min-w-0 flex-1 overflow-auto py-3 text-[13px] leading-relaxed [tab-size:2] [&_pre]:bg-transparent! [&_pre]:[tab-size:2]"
              dangerouslySetInnerHTML={{ __html: todo }}
            />
          </div>
          <div className="flex flex-col gap-6">
            <WindowFrame title={t.code.result}>
              <LiveCanvas name="todo" className="h-[320px]" />
            </WindowFrame>
            <ul className="flex flex-col gap-3">
              {t.code.points.map((point) => (
                <li key={point} className="flex gap-3 text-fd-muted-foreground">
                  <Check className="mt-0.5 size-5 shrink-0 text-fd-primary" />
                  {point}
                </li>
              ))}
            </ul>
          </div>
        </div>
      </Section>

      {/* Features */}
      <Section title={t.features.title} lead={t.features.lead}>
        <div className="grid gap-px overflow-hidden rounded-2xl border bg-fd-border sm:grid-cols-2 lg:grid-cols-3">
          {t.features.items.map((item, i) => {
            const Icon = featureIcons[i]
            return (
              <div key={item.title} className="bg-fd-card p-7">
                <div className="mb-5 inline-flex size-10 items-center justify-center rounded-xl bg-fd-primary/10 text-fd-primary">
                  <Icon className="size-5" />
                </div>
                <h3 className="font-semibold">{item.title}</h3>
                <p className="mt-2 text-sm leading-relaxed text-fd-muted-foreground">
                  {item.body}
                </p>
              </div>
            )
          })}
        </div>
      </Section>

      {/* Native */}
      <Section title={t.native.title} lead={t.native.lead}>
        <img
          src={galleryUrl}
          alt={t.native.alt}
          loading="lazy"
          className="mx-auto w-full max-w-5xl rounded-2xl border shadow-[0_24px_80px_-24px_rgb(0_0_0/0.3)]"
        />
      </Section>

      {/* Components */}
      <Section title={t.components.title} lead={t.components.lead}>
        <div className="flex flex-wrap justify-center gap-2">
          {components.map((c) => (
            <Link
              key={c.url}
              to={c.url}
              prefetch="intent"
              className="rounded-full border bg-fd-card px-3.5 py-1.5 text-sm transition-colors hover:border-fd-primary/50 hover:text-fd-primary"
            >
              {c.title}
            </Link>
          ))}
        </div>
        <div className="mt-10 flex justify-center">
          <Link to={allComponents} className={secondaryLink}>
            {t.components.all}
            <ArrowRight className="size-4" />
          </Link>
        </div>
      </Section>

      {/* Call to action */}
      <section className="border-t">
        <div className="mx-auto flex max-w-6xl flex-col items-center px-6 py-24 text-center">
          <h2 className="text-3xl font-semibold tracking-tight sm:text-4xl">
            {t.cta.title}
          </h2>
          <p className="mt-4 text-lg text-fd-muted-foreground">{t.cta.lead}</p>
          <div className="mt-8 flex flex-wrap justify-center gap-3">
            <Link to={quickStart} className={primaryLink}>
              {t.cta.start}
              <ArrowRight className="size-4" />
            </Link>
            <a
              href={site.github}
              target="_blank"
              rel="noreferrer"
              className={secondaryLink}
            >
              <GithubIcon className="size-4" />
              {t.cta.github}
            </a>
          </div>
        </div>
      </section>

      <footer className="border-t">
        <div className="mx-auto flex max-w-6xl flex-wrap items-center justify-between gap-4 px-6 py-8 text-sm text-fd-muted-foreground">
          <p>© ONEO · {t.footer.license}</p>
          <a
            href={site.mygo}
            target="_blank"
            rel="noreferrer"
            className="hover:text-fd-foreground"
          >
            {t.footer.mygo}
          </a>
        </div>
      </footer>
    </HomeLayout>
  )
}
