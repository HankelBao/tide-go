package main

import (
	"github.com/gdamore/tcell"
)

type TextEditor struct {
	displayRange *DisplayRange
	textBuffer   *TextBuffer
}

func NewTextEditor(firstTextBuffer *TextBuffer) *TextEditor {
	te := new(TextEditor)
	te.displayRange = new(DisplayRange)
	firstTextBuffer.UpdateDisplayRange(te.displayRange)
	te.textBuffer = firstTextBuffer
	return te
}

func (te *TextEditor) SwitchTextBuffer(newTB *TextBuffer) {
	te.textBuffer = newTB
	newTB.UpdateDisplayRange(te.displayRange)
	te.Display()
}

func (te *TextEditor) Display() {
	displayLines := te.textBuffer.GetLines()
	te.displayRange.Display(displayLines, te.textBuffer.lineOffset)
	te.displayRange.ShowCursor(te.textBuffer.GetCursorInDisplayRange())
	screen.Show()
}

func (te *TextEditor) Key(eventKey *tcell.EventKey) {
	key := eventKey.Key()
	switch key {
	case tcell.KeyRune:
		var inputRune = eventKey.Rune()
		te.textBuffer.Insert(inputRune)
	case tcell.KeyUp:
		te.textBuffer.CursorUp()
	case tcell.KeyDown:
		te.textBuffer.CursorDown()
	case tcell.KeyLeft:
		te.textBuffer.CursorLeft()
	case tcell.KeyRight:
		te.textBuffer.CursorRight()
	case tcell.KeyBackspace2:
		te.textBuffer.Backspace()
	case tcell.KeyEnter:
		te.textBuffer.Return()
	case tcell.KeyTab:
		te.textBuffer.Tab()
	case tcell.KeyCtrlS:
		te.textBuffer.Save()
	case tcell.KeyCtrlQ:
		signalDoExit <- true
	}
	te.Display()
}
