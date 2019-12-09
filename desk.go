package desktop // import "fyne.io/desktop"

import (
	"fmt"
	"log"
	"math"
	"runtime/debug"

	"fyne.io/fyne"
)

// Desktop defines an embedded or full desktop environment that we can run.
type Desktop interface {
	Root() fyne.Window
	Run()
	RunApp(AppData) error
	Settings() DeskSettings
	ContentSizePixels(screen *Head) (uint32, uint32)

	IconProvider() ApplicationProvider
	WindowManager() WindowManager
}

var instance Desktop

type deskLayout struct {
	app      fyne.App
	win      fyne.Window
	wm       WindowManager
	icons    ApplicationProvider
	screens  Screens
	settings DeskSettings

	backgrounds         []fyne.CanvasObject
	bar, widgets, mouse fyne.CanvasObject
	container           *fyne.Container
	screenBackgroundMap map[*Head]fyne.CanvasObject
}

func (l *deskLayout) Layout(objs []fyne.CanvasObject, size fyne.Size) {
	var x, y, w, h int = 0, 0, 0, 0
	screens := l.screens.Screens()
	primary := l.screens.Primary()
	if l.Root().Canvas().Scale() != l.screens.Scale() {
		l.Root().Canvas().SetScale(l.screens.Scale())
	}
	x, y, w, h = primary.ScaledX, primary.ScaledY,
		primary.ScaledWidth, primary.ScaledHeight
	size.Width = w
	size.Height = h
	if screens != nil && len(screens) > 1 && len(l.backgrounds) > 1 {
		for i := 0; i < len(screens); i++ {
			if screens[i] == primary {
				continue
			}
			xx, yy, ww, hh := screens[i].ScaledX, screens[i].ScaledY, screens[i].ScaledWidth, screens[i].ScaledHeight
			background := l.screenBackgroundMap[screens[i]]
			if background != nil {
				background.Move(fyne.NewPos(xx, yy))
				background.Resize(fyne.NewSize(ww, hh))
			}
		}
	}

	barHeight := l.bar.MinSize().Height
	l.bar.Resize(fyne.NewSize(size.Width, y+barHeight))
	l.bar.Move(fyne.NewPos(x, y+size.Height-barHeight))

	widgetsWidth := l.widgets.MinSize().Width
	l.widgets.Resize(fyne.NewSize(widgetsWidth, size.Height))
	l.widgets.Move(fyne.NewPos(x+size.Width-widgetsWidth, y))

	background := l.screenBackgroundMap[primary]
	if background != nil {
		background.Move(fyne.NewPos(x, y))
		background.Resize(size)
	}
}

func (l *deskLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(1280, 720)
}

func (l *deskLayout) newDesktopWindow() fyne.Window {
	desk := l.app.NewWindow("Fyne Desktop")
	desk.SetPadded(false)

	if l.wm != nil {
		desk.FullScreen()
	}

	return desk
}

func (l *deskLayout) Root() fyne.Window {
	if l.win == nil {
		l.win = l.newDesktopWindow()

		l.backgrounds = append(l.backgrounds, newBackground())
		l.screenBackgroundMap[l.screens.Screens()[0]] = l.backgrounds[0]
		l.bar = newBar(l)
		l.widgets = newWidgetPanel(l)
		l.mouse = newMouse()
		l.container = fyne.NewContainerWithLayout(l, l.backgrounds[0])
		if l.screens.Screens() != nil && len(l.screens.Screens()) > 1 {
			for i := 1; i < len(l.screens.Screens()); i++ {
				l.backgrounds = append(l.backgrounds, newBackground())
				l.screenBackgroundMap[l.screens.Screens()[i]] = l.backgrounds[i]
				l.container.AddObject(l.backgrounds[i])
			}
		}
		l.container.AddObject(l.bar)
		l.container.AddObject(l.widgets)
		l.container.AddObject(mouse)

		l.win.SetContent(l.container)
		l.mouse.Hide() // temporarily we do not handle mouse (using X default)
		if l.wm != nil {
			l.win.SetOnClosed(func() {
				l.wm.Close()
			})
		}
	}

	return l.win
}

func (l *deskLayout) Run() {
	debug.SetPanicOnFault(true)

	defer func() {
		if r := recover(); r != nil {
			log.Println("Crashed!!!")
			if l.wm != nil {
				l.wm.Close() // attempt to close cleanly to leave X server running
			}
		}
	}()

	l.Root().ShowAndRun()
}

func (l *deskLayout) RunApp(app AppData) error {
	vars := l.scaleVars(l.Root().Canvas().Scale())
	return app.Run(vars)
}

func (l *deskLayout) Settings() DeskSettings {
	return l.settings
}

func (l *deskLayout) ContentSizePixels(head *Head) (uint32, uint32) {
	screenW := uint32(head.Width)
	screenH := uint32(head.Height)
	if l.screens.Primary() == head {
		return screenW - uint32(float32(l.widgets.Size().Width)*l.Root().Canvas().Scale()), screenH
	}
	return screenW, screenH
}

func (l *deskLayout) IconProvider() ApplicationProvider {
	return l.icons
}

func (l *deskLayout) WindowManager() WindowManager {
	return l.wm
}

func (l *deskLayout) scaleVars(scale float32) []string {
	intScale := int(math.Round(float64(scale)))

	return []string{
		fmt.Sprintf("QT_SCALE_FACTOR=%1.1f", scale),
		fmt.Sprintf("GDK_SCALE=%d", intScale),
		fmt.Sprintf("ELM_SCALE=%1.1f", scale),
	}
}

func (l *deskLayout) Screens() []*Head {
	if l.app == nil {
		return nil
	}
	return []*Head{{"Screen0", 0, 0,
		int(l.Root().Canvas().Size().Width), int(l.Root().Canvas().Size().Height),
		0, 0, int(float32(l.Root().Canvas().Size().Width) * l.Root().Canvas().Scale()),
		int(float32(l.Root().Canvas().Size().Height) * l.Root().Canvas().Scale())}}
}

func (l *deskLayout) Active() *Head {
	return l.Screens()[0]
}

func (l *deskLayout) Primary() *Head {
	return l.Screens()[0]
}

func (l *deskLayout) Scale() float32 {
	return l.Root().Canvas().Scale()
}

func (l *deskLayout) ScreenForWindow(windwoX int, windowY int) *Head {
	return l.Screens()[0]
}

// Instance returns the current desktop environment and provides access to injected functionality.
func Instance() Desktop {
	return instance
}

// NewDesktop creates a new desktop in fullscreen for main usage.
// The WindowManager passed in will be used to manage the screen it is loaded on.
// An ApplicationProvider is used to lookup application icons from the operating system.
func NewDesktop(app fyne.App, wm WindowManager, icons ApplicationProvider, screens Screens) Desktop {
	instance = &deskLayout{app: app, wm: wm, icons: icons, screens: screens, settings: NewDeskSettings()}
	instance.(*deskLayout).screenBackgroundMap = make(map[*Head]fyne.CanvasObject)
	return instance
}

// NewEmbeddedDesktop creates a new windowed desktop for test purposes.
// An ApplicationProvider is used to lookup application icons from the operating system.
// If run during CI for testing it will return an in-memory window using the
// fyne/test package.
func NewEmbeddedDesktop(app fyne.App, icons ApplicationProvider) Desktop {
	instance = &deskLayout{app: app, icons: icons, settings: NewDeskSettings()}
	instance.(*deskLayout).screens = instance.(*deskLayout)
	instance.(*deskLayout).screenBackgroundMap = make(map[*Head]fyne.CanvasObject)
	return instance
}
