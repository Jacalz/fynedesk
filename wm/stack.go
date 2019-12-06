package wm

import (
	"fyne.io/desktop"
	"github.com/BurntSushi/xgb/xproto"
)

type stack struct {
	clients      []desktop.Window
	mappingOrder []desktop.Window

	listeners []desktop.StackListener
}

func (s *stack) addToStack(win desktop.Window) {
	s.clients = append([]desktop.Window{win}, s.clients...)
	s.mappingOrder = append(s.mappingOrder, win)
}

func (s *stack) addToStackBottom(win desktop.Window) {
	s.clients = append(s.clients, win)
	s.mappingOrder = append(s.mappingOrder, win)
}

func (s *stack) removeFromStack(win desktop.Window) {
	pos := -1
	for i, w := range s.clients {
		if w == win {
			pos = i
		}
	}

	if pos == -1 {
		return
	}
	s.clients = append(s.clients[:pos], s.clients[pos+1:]...)

	pos = -1
	for i, w := range s.mappingOrder {
		if w == win {
			pos = i
		}
	}
	if pos == -1 {
		return
	}
	s.mappingOrder = append(s.mappingOrder[:pos], s.mappingOrder[pos+1:]...)
}

func (s *stack) getMappingOrder() []xproto.Window {
	return s.getWindowsFromClients(s.mappingOrder)
}

func (s *stack) getWindowsFromClients(clients []desktop.Window) []xproto.Window {
	var wins []xproto.Window
	for _, cli := range clients {
		wins = append(wins, cli.(*client).id)
	}
	return wins
}

func (s *stack) AddWindow(win desktop.Window) {
	if win == nil {
		return
	}
	s.addToStack(win)

	for _, l := range s.listeners {
		l.WindowAdded(win)
	}
}

func (s *stack) RemoveWindow(win desktop.Window) {
	s.removeFromStack(win)

	if s.TopWindow() != nil {
		s.TopWindow().Focus()
	}

	for _, l := range s.listeners {
		l.WindowRemoved(win)
	}
}

func (s *stack) TopWindow() desktop.Window {
	if len(s.clients) == 0 {
		return nil
	}
	return s.clients[0]
}

func (s *stack) Windows() []desktop.Window {
	return s.clients
}

func (s *stack) RaiseToTop(win desktop.Window) {
	if win.Iconic() {
		return
	}
	if len(s.clients) > 1 {
		win.RaiseAbove(s.TopWindow())
	}

	s.removeFromStack(win)
	s.addToStack(win)
}
