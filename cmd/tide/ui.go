package main

import (
	"github.com/gdamore/tcell"
)

type UIElement interface {
	Display()
}

type UISelector interface {
	Key(*tcell.EventKey)
	Focus()
	UnFocus()
}

func RefreshAllUIElements() {
	textEditor.Display()
	statusLine.Display()
	fuzzySwitcher.Display()
}
