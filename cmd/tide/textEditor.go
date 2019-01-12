package main

import (
	"github.com/gdamore/tcell"
)

type TextEditor struct {
	displayRange *DisplayRange
	onFocus      bool

	textBuffer *TextBuffer
}

func InitTextEditor(firstTextBuffer *TextBuffer) *TextEditor {
	te := new(TextEditor)
	te.displayRange = NewDisplayRange()
	te.onFocus = false
	firstTextBuffer.UpdateDisplayRange(te.displayRange)
	te.textBuffer = firstTextBuffer
	return te
}

func (te *TextEditor) OpenFile(fileName string) {
	if fileName == "" {
		return
	}
	var targetTextBuffer *TextBuffer = nil
	for _, textBuffer := range textBuffers {
		if textBuffer.url == fileName {
			targetTextBuffer = textBuffer
			break
		}
	}
	if targetTextBuffer == nil {
		targetTextBuffer = NewTextBuffer()
		targetTextBuffer.Load(fileName)
		textBuffers = append(textBuffers, targetTextBuffer)
	}
	te.SwitchTextBuffer(targetTextBuffer)
}

func (te *TextEditor) SwitchTextBuffer(newTB *TextBuffer) {
	te.textBuffer = newTB
	newTB.UpdateDisplayRange(te.displayRange)
	te.Display()
}

func (te *TextEditor) Display() {
	displayLines := te.textBuffer.GetLines()
	te.displayRange.Display(displayLines, te.textBuffer.lineOffset)
	if te.onFocus == true {
		te.displayRange.ShowCursor(te.textBuffer.GetCursorInDisplayRange())
	}
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
	}
	textEditor.Display()
	statusLine.Display()
}

func (te *TextEditor) Focus() {
	te.onFocus = true
	te.Display()
}

func (te *TextEditor) UnFocus() {
	te.onFocus = false
}
