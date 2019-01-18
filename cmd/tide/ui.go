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
	fileSelector.Display()
}

func InitUIFocus(uiSelector UISelector) {
	focusUISelector = uiSelector
	focusUISelector.Focus()
}

func SwitchUIFocus(targetUISelector UISelector) {
	if focusUISelector == targetUISelector {
		return
	}
	focusUISelector.UnFocus()
	focusUISelector = UISelector(targetUISelector)
	focusUISelector.Focus()
}
