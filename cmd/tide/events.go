package main

import (
	"github.com/gdamore/tcell"
)

// EventLoop is the Main Event Loop
func EventLoop() {
	for eventLoopOn {
		var event tcell.Event
		event = screen.PollEvent()
		switch e := event.(type) {
		case *tcell.EventKey:
			if focusUISelector != nil {
				focusUISelector.Key(e)
			}
		case *tcell.EventResize:
			DisplayUIElements()
		}
	}
}
