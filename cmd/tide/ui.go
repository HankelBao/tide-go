package main

import (
	"github.com/gdamore/tcell"
)

type UIElement interface {
	Display()
}

type UISelector interface {
	Key(*tcell.EventKey)
}

func DisplayUIElements() {
	for _, uiElement := range UIElements {
		uiElement.Display()
	}
}
