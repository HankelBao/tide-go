package main

const HintLineWidth = 30

type HintLine struct {
	displayRange *DisplayRange
	visible      bool

	line    *Line
	content string
}

func InitHintLine() *HintLine {
	hintLine := new(HintLine)
	hintLine.displayRange = NewDisplayRange()
	hintLine.displayRange.width.AbsoluteLength(HintLineWidth)
	hintLine.displayRange.height.AbsoluteLength(1)
	hintLine.line = NewLine()
	hintLine.line.lineStyle.WithDefault(*colorTheme.reversedStyle)
	hintLine.visible = false
	return hintLine
}

func (hl *HintLine) Display() {
	if hl.visible == false {
		return
	}
	xPos, yPos := textEditor.displayRange.GetScreenPos(textEditor.textBuffer.cursor.Get())
	hl.displayRange.horizentalOffset.AbsoluteLength(xPos)
	hl.displayRange.verticalOffset.AbsoluteLength(yPos - 1)
	hl.line.Load(hl.content)
	hl.displayRange.Display([]*Line{hl.line}, 0)
	screen.Show()
}

func (hl *HintLine) LoadContent(content string) {
	hl.content = content
	hl.visible = true
}
