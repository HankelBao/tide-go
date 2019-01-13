package main

type InputLine struct {
	title string
	line  *Line

	cursorStart  int
	cursorOffset int
}

func NewInputLine(title string) *InputLine {
	il := new(InputLine)
	il.title = title
	il.line = NewLine()
	il.UpdateTitle(title)
	return il
}

func (il *InputLine) UpdateTitle(title string) {
	il.line.Load(title)
	il.cursorStart = len(title)
	il.cursorOffset = len(title)
}

func (il *InputLine) CursorLeft() {
	if il.cursorOffset > il.cursorStart {
		il.cursorOffset--
	}
}

func (il *InputLine) CursorRight() {
	if il.cursorOffset < il.line.Len() {
		il.cursorOffset++
	}
}

func (il *InputLine) Insert(r rune) {
	il.line.Insert(il.cursorOffset, byte(r))
	il.CursorRight()
}

func (il *InputLine) Backspace() {
	if il.cursorOffset <= il.cursorStart {
		return
	}
	il.line.Delete(il.cursorOffset)
	il.CursorLeft()
}

func (il *InputLine) GetLine() *Line {
	return il.line
}

func (il *InputLine) GetInput() string {
	title_length := len(il.title)
	return il.line.String()[title_length:]
}
