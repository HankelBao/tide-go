package main

import "bytes"
import "github.com/alecthomas/chroma"

const GapSize = 10

type Line struct {
	left      []byte
	gapBuffer []byte
	right     []byte

	lineStyle *LineStyle
}

func NewLine() *Line {
	line := new(Line)
	line.gapBuffer = make([]byte, 0, GapSize)
	line.lineStyle = NewLineStyle()
	return line
}

func (line *Line) Load(content string) {
	line.left = []byte(content)
	line.gapBuffer = line.gapBuffer[:0]
	line.right = nil
}

func (line *Line) Insert(offset int, content byte) {
	if offset == len(line.left)+len(line.gapBuffer) {
		// Offset is at the right position
		// Only Insert into the Gap Buffer when its len is smaller than GapSize
		// Otherwise the Gap Buffer will be reallocated which ruins the design
		if len(line.gapBuffer) < GapSize {
			line.gapBuffer = append(line.gapBuffer, content)
		} else {
			line.left = append(line.left, line.gapBuffer...)
			line.gapBuffer = line.gapBuffer[:0]
			line.gapBuffer = append(line.gapBuffer, content)
		}
	} else {
		original := line.Bytes()
		line.left = original[:offset]
		line.right = original[offset:]
		line.gapBuffer = line.gapBuffer[:0]
		line.gapBuffer = append(line.gapBuffer, content)
	}
}

func (line *Line) Backspace(offset int) {
	if offset == len(line.left)+len(line.gapBuffer) {
		if len(line.gapBuffer) > 0 {
			line.gapBuffer = line.gapBuffer[:len(line.gapBuffer)-1]
		} else {
			line.left = line.left[:len(line.left)-1]
		}
	} else {
		original := line.Bytes()
		line.left = original[:offset]
		line.right = original[offset:]
		if len(line.left) == 0 {
			return
		}
		line.left = line.left[:len(line.left)-1]
		line.gapBuffer = line.gapBuffer[:0]
	}
}

func (line *Line) JoinToEnd(content []byte) {
	line.right = append(line.right, content...)
}

func (line *Line) CutContentToEnd(offset int) []byte {
	origin := line.Bytes()
	line.left = origin[:offset]
	line.gapBuffer = line.gapBuffer[:0]
	line.right = nil
	return origin[offset:]
}

func (line *Line) ContentBuffer() bytes.Buffer {
	var buffer bytes.Buffer
	for _, b := range line.left {
		buffer.WriteByte(b)
	}
	for _, b := range line.gapBuffer {
		buffer.WriteByte(b)
	}
	for _, b := range line.right {
		buffer.WriteByte(b)
	}
	return buffer
}

func (line *Line) Len() int {
	return len(line.left) + len(line.gapBuffer) + len(line.right)
}

func (line *Line) String() string {
	contentBuffer := line.ContentBuffer()
	return contentBuffer.String()
}

func (line *Line) Bytes() []byte {
	contentBuffer := line.ContentBuffer()
	return contentBuffer.Bytes()
}

func (line *Line) Debug() string {
	var buffer bytes.Buffer
	for _, b := range line.left {
		buffer.WriteByte(b)
	}
	for _, b := range line.gapBuffer {
		buffer.WriteByte(b)
	}
	for i := 0; i < GapSize-len(line.gapBuffer); i++ {
		buffer.WriteByte('~')
	}
	for _, b := range line.right {
		buffer.WriteByte(b)
	}
	return buffer.String()
}

func (line *Line) UpdateLineStyle(lexer chroma.Lexer) {
	if lexer == nil {
		return
	}
	line.lineStyle.LoadFromSourceCode(line.String(), lexer)
}
