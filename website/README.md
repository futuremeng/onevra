# Onevra website

The documentation site: guides and a page for every component, in English and Simplified Chinese. Every example on it is Go code from [`source/examples/website`](../source/examples/website), compiled to WebAssembly and running in the page. It deploys to GitHub Pages at <https://oneo-me.github.io/onevra/>.

The design and its constraints are in [Website design](../docs/architecture/website.md).

## Requirements

- [Bun](https://bun.sh) 1.4 or newer.
- [Go](https://go.dev) 1.27.3 or newer, which compiles the examples and reads the API.

## Commands

Run these from `website/`:

```sh
bun install          # install dependencies
bun run dev          # serve at http://localhost:5173/onevra/
bun run check        # check that the docs cover the API, examples and translations
bun run typecheck    # generate route types and typecheck
bun run build        # prerender the static site into build/pages
bun run preview      # serve build/pages the way GitHub Pages does
bun run format       # format with Prettier
```

The development server rebuilds the examples and recompiles the pages when Go source changes.

## Adding an example

1. Add a file to `source/examples/website/examples/`, such as `slider_range.go`. Its name is the example's: `slider-range`.
2. Declare the example's state as a type with a `View(c *ui.Context)` method, and register one value with `var _ = register(&value{})`. Use `registerWindow` when it lays out the whole window itself.
3. Show it on a page with `<Example name="slider-range" height={240} />`, in both the English and the Chinese page.

The page shows the file's source without its package clause and registration. Write it as an application would, with plain English labels; `go test ./examples/website/...` from `source/` builds every example.

## Adding a page

1. Add `name.mdx` and `name.zh.mdx` under `content/docs/components/` or `content/docs/(guides)/`, each with a `title` and a `description`.
2. List the page in that folder's `meta.json` and `meta.zh.json`, which order the sidebar.
3. End a component page with an `## API` section of `<GoAPI name="components.Name" />` entries, which read the declaration and doc comment from Go.

Pages can also use Fumadocs' `Steps`, `Tabs`, `Cards`, `Callout` and `TypeTable`, registered in [`app/components/mdx.tsx`](app/components/mdx.tsx).

## Structure

| Path | Purpose |
| --- | --- |
| `content/docs/` | Guides and component pages, with their translations and sidebar order |
| `app/routes/` | The home page, the docs pages and the static search index |
| `app/components/home/` | The home page and its copy in both languages |
| `app/components/example.tsx` | A live example above its source |
| `app/lib/examples.ts` | The browser side of the examples' WebAssembly |
| `scripts/examples.ts` | Vite plugin that writes `.generated/docs-data.json` and compiles the WebAssembly |
| `scripts/remark-onevra.ts` | Fills `<Example>` and `<GoAPI>` from that data while compiling MDX |
| `scripts/check-docs.ts` | The coverage check behind `bun run check` |

`app/assets/icon.png` is downscaled from `../source/assets/icon.png`. Regenerate it rather than editing it by hand. Add shadcn/ui components with `bun x --bun shadcn@latest add <component>`.

## Deployment

The [Website workflow](../.github/workflows/website.yml) checks, builds and publishes the site whenever `website/**` or `source/**` changes on `main`. Enable it once under **Settings → Pages → Build and deployment → Source → GitHub Actions**.

`BASE_PATH` overrides the deploy path in `vite.config.ts` and `react-router.config.ts`. It defaults to `/onevra/`, which matches the `oneo-me/onevra` repository; set it to `/` for a custom domain. React Router writes prerendered pages under `build/client/onevra/` while Vite keeps assets at the root of `build/client`, so `scripts/prepare-pages.mjs` merges both into `build/pages` and adds a `404.html` fallback.
