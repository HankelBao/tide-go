package main

import (
	"github.com/gdamore/tcell"
)

func ExecuteGlobalBinding(eventKey *tcell.EventKey) bool {
	key := eventKey.Key()
	switch key {
	case tcell.KeyCtrlQ:
		signalDoExit <- true
	case tcell.KeyEsc:
		SwitchUIFocus(textEditor)
	case tcell.KeyCtrlP:
		SwitchUIFocus(fileSelector)
	default:
		return false
	}
	return true
}
