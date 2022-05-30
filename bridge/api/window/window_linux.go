package window

import (
  "log"

  "tractor.dev/apptron/bridge/resource"
  "tractor.dev/apptron/bridge/platform/linux"
)

type Window struct {
  window

  Window linux.Window
  Webview linux.Webview

  callbackId int
}

func init() {
  linux.OS_Init()
}

func New(options Options) (*Window, error) {
  win := &Window{
    window: window{
      Handle: resource.NewHandle(),
    },
  }
  resource.Retain(win.Handle, win)


  window := linux.Window_New()

  window.SetPosition(int(options.Position.X), int(options.Position.Y))
  window.SetSize(int(options.Size.Width), int(options.Size.Height))

  if options.MinSize.Width != 0 || options.MinSize.Height != 0 {
    window.SetMinSize(int(options.MinSize.Width), int(options.MinSize.Height))
  }

  if options.MaxSize.Width != 0 || options.MaxSize.Height != 0 {
    window.SetMaxSize(int(options.MaxSize.Width), int(options.MaxSize.Height))
  }

  if options.Center {
    window.Center()
  }

  if options.Frameless {
    window.SetDecorated(false)
  }

  if options.Fullscreen {
    window.SetFullscreen(true)
  }

  if options.Maximized {
    window.SetMaximized(true)
  }

  window.SetResizable(options.Resizable)

  if options.Title != "" {
    window.SetTitle(options.Title)
  }

  if options.AlwaysOnTop {
    window.SetAlwaysOnTop(true)
  }

  if len(options.Icon) > 0 {
    window.SetIconFromBytes(options.Icon)
  }

  webview := linux.Webview_New()
  webview.SetSettings(linux.DefaultWebviewSettings())

  myCallback := func(result string) {
    log.Println("Callback from JavaScript!!", result)
  }
  callbackId := webview.RegisterCallback("apptron", myCallback)
  webview.Eval("webkit.messageHandlers.apptron.postMessage(JSON.stringify({ hello: 42 }));")

  window.AddWebview(webview)

  if options.Transparent {
    window.SetTransparent(true)
    webview.SetTransparent(true)
  }

  if options.URL != "" {
    webview.Navigate(options.URL)
  }

  if options.HTML != "" {
    webview.SetHtml(options.HTML)
  }

  if options.Script != "" {
    webview.AddScript(options.Script)
  }

  if options.Visible {
    window.Show()
  }

  win.Window  = window
  win.Webview = webview
  win.callbackId = callbackId

  return win, nil
}

func (w *Window) Destroy() {
  if w.callbackId != 0 {
    linux.UnregisterCallback(w.callbackId)
    w.callbackId = 0
  }

  w.Webview.Destroy()
  w.Window.Destroy()
}

func (w *Window) Focus() {
  w.Window.Focus()
}

func (w *Window) SetVisible(visible bool) {
  if visible {
    w.Window.Show()
  } else {
    w.Window.Hide()
  }
}

func (w *Window) IsVisible() bool {
  return w.Window.IsVisible()
}

func (w *Window) SetMaximized(maximized bool) {
  w.Window.SetMaximized(maximized)
}

func (w *Window) SetMinimized(minimized bool) {
  w.Window.SetMinimized(minimized)
}

func (w *Window) SetFullscreen(fullscreen bool) {
  w.Window.SetFullscreen(fullscreen)
}

func (w *Window) SetSize(size Size) {
 w.Window.SetSize(int(size.Width), int(size.Height))
}

func (w *Window) SetMinSize(size Size) {
  w.Window.SetMinSize(int(size.Width), int(size.Height))
}

func (w *Window) SetMaxSize(size Size) {
  w.Window.SetMaxSize(int(size.Width), int(size.Height))
}

func (w *Window) SetResizable(resizable bool) {
  w.Window.SetResizable(resizable)
}

func (w *Window) SetAlwaysOnTop(always bool) {
  w.Window.SetAlwaysOnTop(always)
}

func (w *Window) SetPosition(position Position) {
  w.Window.SetPosition(int(position.X), int(position.Y))
}

func (w *Window) SetTitle(title string) {
  w.Window.SetTitle(title)
}

func (w *Window) GetOuterPosition() Position {
  pos := w.Window.GetPosition()
  return Position{
    X: float64(pos.X),
    Y: float64(pos.Y),
  }
}

func (w *Window) GetOuterSize() Size {
  size := w.Window.GetSize()
  return Size{
    Width:  float64(size.Width),
    Height: float64(size.Height),
  }
}
