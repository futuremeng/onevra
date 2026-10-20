// The browser side of the live examples: one WebAssembly instance of
// source/examples/website for the whole site, which opens a MyGo window for
// each example on the page. MyGo puts a new window in the element with the ID
// "mygo", so mounts take turns giving that ID to their element.
import { loaderUrl, wasmUrl } from "virtual:onevra-examples"

interface ExamplesAPI {
  mount(name: string): Promise<number>
  unmount(handle: number): void
  setLocale(tag: string): void
}

declare global {
  interface Window {
    onevraExamples?: ExamplesAPI
    Go?: new () => {
      argv: string[]
      importObject: WebAssembly.Imports
      run(instance: WebAssembly.Instance): Promise<void>
    }
  }
}

/** Each window has a WebGL context, and browsers keep only a few. */
const maxWindows = 8

let starting: Promise<ExamplesAPI> | undefined
let queue: Promise<unknown> = Promise.resolve()
let language = ""

function mygoLanguage(): string {
  return document.documentElement.lang.startsWith("zh") ? "zh-Hans" : "en"
}

function loadScript(src: string): Promise<void> {
  return new Promise((resolve, reject) => {
    const script = document.createElement("script")
    script.src = src
    script.onload = () => resolve()
    script.onerror = () => reject(new Error(`could not load ${src}`))
    document.head.append(script)
  })
}

async function instantiate(imports: WebAssembly.Imports) {
  try {
    return await WebAssembly.instantiateStreaming(fetch(wasmUrl), imports)
  } catch {
    // Some static hosts serve .wasm as another type; compile the bytes
    // rather than failing on the content type.
    const response = await fetch(wasmUrl)
    return WebAssembly.instantiate(await response.arrayBuffer(), imports)
  }
}

/**
 * MyGo draws in the light or dark appearance that data-mygo-theme on the
 * root element names; keep it on the site's theme.
 */
function followTheme() {
  const root = document.documentElement
  const sync = () =>
    root.setAttribute(
      "data-mygo-theme",
      root.classList.contains("dark") ? "dark" : "light"
    )
  sync()
  new MutationObserver(sync).observe(root, {
    attributes: true,
    attributeFilter: ["class"],
  })
}

function start(): Promise<ExamplesAPI> {
  starting ??= (async () => {
    followTheme()
    if (!window.Go) await loadScript(loaderUrl)
    const go = new window.Go!()
    language = mygoLanguage()
    go.argv = ["onevra-examples", "-lang", language]
    const ready = new Promise<ExamplesAPI>((resolve) =>
      window.addEventListener(
        "onevra-examples-ready",
        () => resolve(window.onevraExamples!),
        { once: true }
      )
    )
    const { instance } = await instantiate(go.importObject)
    void go.run(instance)
    return ready
  })().catch((error) => {
    starting = undefined
    throw error
  })
  return starting
}

interface Live {
  element: HTMLElement
  handle: number
  visible: boolean
  seen: number
}

const live = new Set<Live>()

/** Closes the windows least recently seen offscreen beyond the limit. */
function evict() {
  const hidden = [...live]
    .filter((w) => !w.visible)
    .sort((a, b) => a.seen - b.seen)
  while (live.size > maxWindows && hidden.length > 0) {
    const w = hidden.shift()!
    live.delete(w)
    window.onevraExamples?.unmount(w.handle)
    w.element.dispatchEvent(new Event("onevra-example-evicted"))
  }
}

/**
 * Opens an example in an element and returns its handle. Mounts run one
 * after another, as each borrows the "mygo" ID.
 */
export function mountExample(
  element: HTMLElement,
  name: string
): Promise<Live> {
  const turn = queue.then(async () => {
    const api = await start()
    const tag = mygoLanguage()
    if (tag !== language) {
      language = tag
      api.setLocale(tag)
    }
    document.getElementById("mygo")?.removeAttribute("id")
    element.id = "mygo"
    // A window's title becomes the page's; keep the page's own.
    const title = document.title
    try {
      const handle = await api.mount(name)
      const w: Live = {
        element,
        handle,
        visible: true,
        seen: performance.now(),
      }
      live.add(w)
      evict()
      return w
    } finally {
      element.removeAttribute("id")
      document.title = title
    }
  })
  queue = turn.catch(() => undefined)
  return turn
}

export function unmountExample(w: Live) {
  if (live.delete(w)) window.onevraExamples?.unmount(w.handle)
}

export function setExampleVisible(w: Live, visible: boolean) {
  w.visible = visible
  w.seen = performance.now()
  if (!visible) evict()
}

/** Starts loading the WebAssembly before an example asks for it. */
export function preloadExamples() {
  void start().catch(() => undefined)
}
