package main

import (
	"github.com/gdamore/tcell"  
)

var (
	screen tcell.Screen
)

func ScreenInit() {
	screen, _ = tcell.NewScreen()
	screen.Init()
}

func ScreenToFini() {
	if screen != nil {
		screen.Fini()
	}
}