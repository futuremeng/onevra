declare module "virtual:onevra-examples" {
  /** The examples' WebAssembly, named by its content in a build. */
  export const wasmUrl: string
  /** Go's loader of WebAssembly, matching the Go that built it. */
  export const loaderUrl: string
}
