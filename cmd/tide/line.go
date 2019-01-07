package main

import "bytes"

const GAP_SIZE = 10

type Line struct {
	left []byte
	gapBuffer []byte
	right []byte
}

func NewLine() *Line {
	line := new(Line)
	line.gapBuffer = make([]byte, 0, GAP_SIZE)
	return line
}

func (line *Line) Load(content string) {
	line.left = []byte(content)
}

func (line *Line) Insert(offset int, content byte) {
	if offset == len(line.left)+len(line.gapBuffer) {
		// Offset is at the right position
		if len(line.gapBuffer) < GAP_SIZE {
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
	for i:=0; i<GAP_SIZE-len(line.gapBuffer); i++ {
		buffer.WriteByte('~')
	}
	for _, b := range line.right {
		buffer.WriteByte(b)
	}
	return buffer.String()
}