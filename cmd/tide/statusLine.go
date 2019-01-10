package main

import "strconv"

type StatusLine struct {
	displayRange *DisplayRange
	statusLine   *Line
}

func InitStatusLine() *StatusLine {
	statusLine := new(StatusLine)
	statusLine.displayRange = NewDisplayRange()
	statusLine.statusLine = NewLine()
	statusLine.statusLine.lineStyle.WithDefault(*colorTheme.reversedStyle)
	return statusLine
}

func (sl *StatusLine) Display() {
	statusLineWidth := sl.displayRange.width.GetAbsoluteValue()

	fileNameInfo := textEditor.textBuffer.name

	lineNumStr := strconv.Itoa(textEditor.textBuffer.cursor.lineNum)
	offsetStr := strconv.Itoa(textEditor.textBuffer.cursor.offset)
	cursorInfo := "Ln " + lineNumStr + ", Col " + offsetStr

	statusLineInfo := CombineLine(statusLineWidth, []string{fileNameInfo}, []string{cursorInfo})

	sl.statusLine.Load(statusLineInfo)

	sl.displayRange.Display([]*Line{sl.statusLine}, 0)

}

func CombineLine(width int, left []string, right []string) string {
	divider := "  "
	combinedLine := make([]byte, width)
	var combinedLeft string
	for index, item := range left {
		if index == 0 {
			combinedLeft = item
		} else {
			combinedLeft = combinedLeft + divider + item
		}
	}
	var combinedRight string
	for _, item := range right {
		combinedRight = combinedRight + divider + item
	}
	if len(combinedLeft) > width || len(combinedRight) > width {
		return ""
	}
	copy(combinedLine[:len(combinedLeft)], combinedLeft)
	copy(combinedLine[width-len(combinedRight):], combinedRight)

	return string(combinedLine)
}
