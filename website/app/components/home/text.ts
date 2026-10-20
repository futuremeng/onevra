import type { Locale } from "~/lib/i18n"

export interface HomeText {
  title: string
  description: string
  badge: string
  headline: [string, string]
  lead: string
  start: string
  browse: string
  copy: string
  copied: string
  live: string
  liveNote: string
  code: { title: string; lead: string; points: string[]; result: string }
  features: {
    title: string
    lead: string
    items: { title: string; body: string }[]
  }
  native: { title: string; lead: string; alt: string }
  components: { title: string; lead: string; all: string }
  cta: { title: string; lead: string; start: string; github: string }
  footer: { license: string; mygo: string }
}

export const homeText: Record<Locale, HomeText> = {
  en: {
    title: "Onevra — Native desktop components for Go",
    description:
      "A Go component library for native desktop applications, built on MyGo: themed controls, window shells, Lucide icons and complex views, drawn without a web view.",
    badge: "Built on MyGo",
    headline: ["Small components.", "Real possibilities."],
    lead: "Onevra is a Go component library for native desktop apps. Themed controls, window shells, icons, charts and chat, all drawn by MyGo without a web view.",
    start: "Get started",
    browse: "Browse components",
    copy: "Copy",
    copied: "Copied",
    live: "Overview — Onevra",
    liveNote:
      "Not a screenshot: this window is Go, compiled to WebAssembly and drawn by MyGo in your browser. Click around.",
    code: {
      title: "Write Go. Get a desktop app.",
      lead: "Builders run every frame and read your application's state. There is no markup, no bridge and no JavaScript.",
      points: [
        "State lives in your own values, passed by pointer.",
        "Controls report clicks and edits as they happen.",
        "The same code runs on macOS, Windows and Linux.",
      ],
      result: "todo.go — running",
    },
    features: {
      title: "Everything an application window needs",
      lead: "Controls that refine MyGo's, components it does not have, and the shells that hold them.",
      items: [
        {
          title: "Themed controls",
          body: "Buttons, inputs, selects, sliders, tabs, tables and dialogs that share one consistent look in light and dark.",
        },
        {
          title: "Window shells",
          body: "Plain and sidebar windows with native window controls, a header for actions and resizable navigation.",
        },
        {
          title: "Charts and diagrams",
          body: "Bar, line, area and donut charts, relationship graphs, treemaps and flow diagrams over data you own.",
        },
        {
          title: "Conversation interfaces",
          body: "Message turns, tool activity, folded summaries and a glass composer for chat-based applications.",
        },
        {
          title: "Lucide icons",
          body: "More than 1,800 Lucide icons as named Go accessors, parsed on first use and tinted by the surrounding text.",
        },
        {
          title: "Tested without a window",
          body: "Drive any view headlessly with MyGo's tester: click, type and check results in an ordinary go test.",
        },
      ],
    },
    native: {
      title: "At home on the desktop",
      lead: "In a real window, Onevra keeps the platform's window controls and system font, with glass where the system has it. The Gallery shows every component this way.",
      alt: "The Onevra Gallery on macOS, with its sidebar and welcome page",
    },
    components: {
      title: "A component for every part of the window",
      lead: "Each one has a page with live examples, their source and the API, read from the Go code itself.",
      all: "See all components",
    },
    cta: {
      title: "Build your first window",
      lead: "Install the module, copy the quick start and run it with go run.",
      start: "Read the quick start",
      github: "Star on GitHub",
    },
    footer: { license: "MIT licensed", mygo: "Powered by MyGo" },
  },
  zh: {
    title: "Onevra — Go 原生桌面组件库",
    description:
      "基于 MyGo 的 Go 原生桌面应用组件库：主题化控件、窗口外壳、Lucide 图标和复杂视图，无需 Web 视图即可绘制。",
    badge: "基于 MyGo 构建",
    headline: ["小组件，", "大可能。"],
    lead: "Onevra 是面向原生桌面应用的 Go 组件库。主题化控件、窗口外壳、图标、图表和对话界面，全部由 MyGo 直接绘制，不依赖 Web 视图。",
    start: "快速开始",
    browse: "浏览组件",
    copy: "复制",
    copied: "已复制",
    live: "概览 — Onevra",
    liveNote:
      "这不是截图：这个窗口是 Go 代码，编译为 WebAssembly 后由 MyGo 在浏览器中绘制。试着点一点。",
    code: {
      title: "写 Go，得到桌面应用。",
      lead: "构建函数每帧运行并读取应用状态。没有标记语言，没有桥接层，也没有 JavaScript。",
      points: [
        "状态保存在你自己的值里，通过指针传入。",
        "控件在发生时报告点击和编辑。",
        "同一份代码可运行在 macOS、Windows 和 Linux 上。",
      ],
      result: "todo.go — 运行中",
    },
    features: {
      title: "应用窗口所需的一切",
      lead: "改进 MyGo 自带的控件，补充它没有的组件，再加上承载它们的窗口外壳。",
      items: [
        {
          title: "主题化控件",
          body: "按钮、输入框、选择器、滑块、标签页、表格和对话框，在浅色与深色模式下保持一致的外观。",
        },
        {
          title: "窗口外壳",
          body: "提供原生窗口控件的普通窗口与侧边栏窗口，带有放置操作的标题栏和可调整宽度的导航。",
        },
        {
          title: "图表与示意图",
          body: "柱状图、折线图、面积图和环形图，以及关系图、矩形树图和流程图，数据由应用掌控。",
        },
        {
          title: "对话界面",
          body: "消息轮次、工具活动、折叠摘要和玻璃质感的输入框，适合基于对话的应用。",
        },
        {
          title: "Lucide 图标",
          body: "1800 多个 Lucide 图标都是具名的 Go 访问器，首次使用时解析，并继承周围文字的颜色。",
        },
        {
          title: "无需窗口即可测试",
          body: "使用 MyGo 的测试器以无界面方式驱动任意视图：在普通的 go test 中点击、输入并检查结果。",
        },
      ],
    },
    native: {
      title: "在桌面上浑然天成",
      lead: "在真实窗口中，Onevra 沿用平台的窗口控件和系统字体，并在系统支持时使用玻璃材质。Gallery 以这种方式展示每个组件。",
      alt: "macOS 上的 Onevra Gallery，带有侧边栏和欢迎页",
    },
    components: {
      title: "窗口的每个部分都有对应组件",
      lead: "每个组件都有独立页面，包含实时示例、源码和直接从 Go 代码读取的 API。",
      all: "查看全部组件",
    },
    cta: {
      title: "构建你的第一个窗口",
      lead: "安装模块，复制快速开始中的代码，然后用 go run 运行。",
      start: "阅读快速开始",
      github: "在 GitHub 上加星",
    },
    footer: { license: "MIT 许可", mygo: "由 MyGo 驱动" },
  },
}
