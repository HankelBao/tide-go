package main

import (
	"github.com/gdamore/tcell"
)

func ExecuteGlobalBinding(eventKey *tcell.EventKey) bool {
	key := eventKey.Key()
	switch key {
	case tcell.KeyCtrlQ:
		signalDoExit <- true
	case tcell.KeyCtrlJ:
		focusUISelector.UnFocus()
		focusUISelector = UISelector(textEditor)
		focusUISelector.Focus()
	case tcell.KeyCtrlSpace:
		focusUISelector.UnFocus()
		focusUISelector = UISelector(fuzzySwitcher)
		focusUISelector.Focus()
	default:
		return false
	}
	return true
}
