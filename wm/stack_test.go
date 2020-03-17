package wm

import (
	"testing"

	"fyne.io/fyne"
	"github.com/stretchr/testify/assert"

	"fyne.io/desktop"
)

type dummyWindow struct {
}

func (*dummyWindow) Class() []string {
	return []string{"dummy1", "dummy2"}
}

func (*dummyWindow) Close() {
	// no-op
}

func (*dummyWindow) Command() string {
	return "DummyCommand"
}

func (*dummyWindow) Decorated() bool {
	return true
}

func (*dummyWindow) Focus() {
	// no-op
}

func (*dummyWindow) Focused() bool {
	return false
}

func (*dummyWindow) Fullscreen() {
	// no-op
}

func (*dummyWindow) Fullscreened() bool {
	return false
}

func (*dummyWindow) Icon() fyne.Resource {
	return nil
}

func (*dummyWindow) Iconic() bool {
	return false
}

func (*dummyWindow) Iconify() {
	// no-op
}

func (*dummyWindow) IconName() string {
	return "DummyIcon"
}

func (*dummyWindow) Maximize() {
	// no-op
}

func (*dummyWindow) Maximized() bool {
	return false
}

func (*dummyWindow) SkipTaskbar() bool {
	return false
}

func (*dummyWindow) Title() string {
	return "Dummy"
}

func (*dummyWindow) TopWindow() bool {
	return true
}

func (*dummyWindow) Unfullscreen() {
	// no-op
}

func (*dummyWindow) Uniconify() {
	// no-op
}

func (*dummyWindow) Unmaximize() {
	// no-op
}

func (*dummyWindow) RaiseAbove(desktop.Window) {
	// no-op (this is instructing the window after stack changes)
}

func (*dummyWindow) RaiseToTop() {
	// no-op
}

func TestStack_AddWindow(t *testing.T) {
	stack := &stack{}
	win := &dummyWindow{}

	stack.AddWindow(win)
	assert.Equal(t, 1, len(stack.Windows()))
}

func TestStack_RaiseToTop(t *testing.T) {
	stack := &stack{}
	win1 := &dummyWindow{}
	win2 := &dummyWindow{}

	stack.AddWindow(win1)
	stack.AddWindow(win2)
	assert.Equal(t, win1, stack.TopWindow())

	stack.RaiseToTop(win2)
	assert.Equal(t, win2, stack.TopWindow())
}

func TestStack_RemoveWindow(t *testing.T) {
	stack := &stack{}
	win := &dummyWindow{}

	stack.AddWindow(win)
	stack.RemoveWindow(win)
	assert.Equal(t, 0, len(stack.Windows()))

}
