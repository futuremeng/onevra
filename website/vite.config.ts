import { reactRouter } from "@react-router/dev/vite"
import tailwindcss from "@tailwindcss/vite"
import { fumadocsMdx } from "fumadocs-mdx/vite"
import { defineConfig } from "vite"

import { onevraExamples } from "./scripts/examples.ts"
import { remarkOnevra } from "./scripts/remark-onevra.ts"

// Must match the React Router `basename` in react-router.config.ts so the
// generated asset URLs also resolve under the GitHub Pages project path.
const basePath = process.env.BASE_PATH ?? "/onevra/"

export default defineConfig({
  base: basePath,
  resolve: { tsconfigPaths: true },
  // The build's output and the generated data change without needing a
  // reload; scripts/examples.ts reloads pages when the Go source changes.
  server: { watch: { ignored: ["**/build/**", "**/.generated/**"] } },
  plugins: [
    onevraExamples(),
    fumadocsMdx({
      globalOptions: {
        mdxOptions: {
          remarkPlugins: (plugins) => [remarkOnevra, ...plugins],
          rehypeCodeOptions: {
            themes: { light: "github-light", dark: "github-dark" },
          },
        },
      },
    }),
    tailwindcss(),
    reactRouter(),
  ],
})
