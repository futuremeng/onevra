# Website design

The website in `website/` is Onevra's public documentation: a static React Router site with Fumadocs, whose live examples are Go code compiled to WebAssembly. Like the Gallery, it is a consumer of the public API, never a dependency of the library.

## Stack and output

| Part | Choice |
| --- | --- |
| Framework | React Router in framework mode, prerendered into static files for GitHub Pages |
| Documentation | Fumadocs with Base UI: layout, sidebar tabs, table of contents, i18n and static search |
| Home page | Tailwind CSS on shadcn/ui's theme tokens, which Fumadocs shares, in Onevra's palette |
| Content | MDX in `website/content/docs`, under the `(guides)` and `components` sidebar roots |
| Examples | [`source/examples/website`](../../source/examples/website/main.go), one WebAssembly app for every example |

The build prerenders every page, with `ssr: false`. Development renders on request instead, so a page added while the server runs opens at once. English pages have no URL prefix and Chinese pages live under `/zh`; each translation sits beside its page as `name.zh.mdx`.

## Examples are the source of truth

Each file of package `examples` is one example, registered under its file name: `button_icons.go` is `button-icons`. A page shows it with `<Example name="button-icons" />`. Examples are written as consumer code, with state in a value and a `View` builder, so the source a reader copies is the source that runs.

- `register` centers an example in its window with a margin; `registerWindow` gives it the whole window, as an application's content.
- `go run ./examples/website -manifest` prints every example's shown source, without its package clause and registration, and the exported API of the documented packages, read with `go/doc`.
- The site writes that output to `website/.generated/docs-data.json`. A remark plugin fills `<Example>` with its source and expands `<GoAPI name="components.Switch" />` into a heading, the declaration, the doc comment, a field table, values and methods.

Filling these at compile time keeps each page's bundle to its own code, highlights Go with the rest of the page, and fails the build on an unknown name. Property definitions and defaults stay in Go doc comments, which the API sections read rather than restate.

## Browser runtime

One WebAssembly instance serves the whole tab and survives client-side navigation. Each example is a MyGo window, which MyGo places in the element with the ID `mygo`; mounts take turns giving that ID to their element, then restore the page's title.

- An example opens when it nears the viewport, hidden and shown without the keyboard, so scrolling never moves focus.
- Each window holds a WebGL context and browsers keep only a few, so at most eight stay open; the least recently seen offscreen one closes first.
- MyGo takes every wheel event over its canvas, so the page keeps the wheel until an example is clicked.
- The appearance follows the site's theme through `data-mygo-theme`, and the language follows the page.

Browsers have no trays, notifications, global shortcuts or native menus; MyGo reports them as unavailable or draws HTML menus. Examples needing them belong in the Gallery only.

## Constraints

- Keep examples product-neutral and self-contained: plain English labels, no I/O in builders, and images decoded once at registration.
- Package `examples` shares one namespace, so prefix helper names with their example.
- Every exported function of a documented package needs a `<GoAPI>` section, every example a page, and every page a translation; `bun run check` in `website/` enforces this.
- Pages describe behavior a reader needs; option-by-option detail stays in doc comments.

## Verification

`go test ./examples/website/...` builds every example in both themes and checks the manifest. The site's typecheck, coverage check and build compile the WebAssembly and prerender every page. None of these prove pointer, keyboard or assistive-technology behavior in a browser, which needs its own review.
