package main

type TextBuffer struct {
	lines []*Line

	url string

}

func NewTextBuffer() *TextBuffer {
	textBuffer := new(TextBuffer)
	return textBuffer
}

func (tb *TextBuffer) GetContent(lineNum int, size int) []*string {
	return_content := make([]*string, 0, size)
	for i:=lineNum; i<=lineNum+size; i++ {
		// Don't worry, Go is magical
		str := tb.lines[i].String()
		return_content = append(return_content, &str)
	}
	return return_content
}

