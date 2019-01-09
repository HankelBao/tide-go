package main

type BufferLine struct {
	displayRange *DisplayRange
	displayLine  *Line
}

func NewBufferLine() *BufferLine {
	bl := new(BufferLine)
	bl.displayRange = new(DisplayRange)
	bl.displayLine = NewLine()
	return bl
}

func (bl *BufferLine) Display() {
	var statusline []byte
	statusline = make([]byte, bl.displayRange.width.GetAbsoluteValue())

	var left []byte
	left = append(left, []byte("File: ")...)
	current_file_name := GetFileNameFromTextBuffer(textEditor.textBuffer)
	left = append(left, current_file_name...)

	var right []byte
	for _, textBuffer := range textBuffers {
		buffer_name := GetFileNameFromTextBuffer(textBuffer)
		right = append(right, buffer_name...)
		right = append(right, '|')
	}

	copy(statusline[:len(left)], left)
	copy(statusline[len(statusline)-len(right):], right)

	bl.displayLine.Load(string(statusline))
	bl.displayRange.Display([]*Line{bl.displayLine}, 0)
}

func GetFileNameFromTextBuffer(tb *TextBuffer) []byte {
	if tb.url == "" {
		return []byte("[No Name]")
	} else {
		return []byte(tb.url)
	}
}
