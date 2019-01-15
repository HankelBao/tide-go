package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/lexers"
)

const TabSize = 4

type TextBuffer struct {
	lines  []*Line
	cursor *Cursor

	topLine    int
	lineOffset int

	displayRange *DisplayRange

	name     string
	url      string
	fileType string
	lexer    chroma.Lexer
}

func NewTextBuffer() *TextBuffer {
	textBuffer := new(TextBuffer)
	textBuffer.cursor = new(Cursor)
	textBuffer.InitFirstLine()
	textBuffer.name = "[No Name]"
	return textBuffer
}

func (tb *TextBuffer) InitFirstLine() {
	firstLine := NewLine()
	tb.lines = []*Line{firstLine}
}

func (tb *TextBuffer) UpdateDisplayRange(dr *DisplayRange) {
	tb.displayRange = dr
}

func (tb *TextBuffer) GetLines() []*Line {
	height := tb.displayRange.height.GetAbsoluteValue()
	returnContent := make([]*Line, 0, height)
	for i := tb.topLine; i <= tb.topLine+height; i++ {
		var line *Line
		if i > len(tb.lines)-1 {
			line = NewLine()
			line.Load("~")
		} else {
			line = tb.lines[i]
		}
		returnContent = append(returnContent, line)
	}
	return returnContent
}

func (tb *TextBuffer) GetCursorInDisplayRange() (int, int) {
	offset, lineNum := tb.cursor.Get()
	return offset - tb.lineOffset, lineNum - tb.topLine
}

func (tb *TextBuffer) CursorLeft() {
	if tb.cursor.offset == 0 {
		return
	}
	tb.cursor.offset--
	if tb.cursor.offset < tb.lineOffset {
		tb.lineOffset--
	}
}

func (tb *TextBuffer) CursorRight() {
	offset, lineNum := tb.cursor.Get()
	if offset+1 > tb.lines[lineNum].Len() {
		return
	}
	tb.cursor.offset++
	if tb.cursor.offset >= tb.lineOffset+tb.displayRange.width.GetAbsoluteValue() {
		tb.lineOffset++
	}
}

func (tb *TextBuffer) CursorOffsetAdjust() {
	lineNum := tb.cursor.lineNum
	if tb.cursor.offset > tb.lines[lineNum].Len() {
		tb.cursor.offset = tb.lines[lineNum].Len()
	}
	// Make sure cursor is still in the screen
	//width := tb.displayRange.width.GetAbsoluteValue()
	if tb.lineOffset > tb.cursor.offset {
		tb.lineOffset = tb.cursor.offset
	}
}

func (tb *TextBuffer) CursorUp() {
	lineNum := tb.cursor.lineNum
	if lineNum == 0 {
		return
	}
	tb.cursor.lineNum--
	if tb.cursor.lineNum < tb.topLine {
		tb.topLine--
	}
	tb.CursorOffsetAdjust()
}

func (tb *TextBuffer) CursorDown() {
	lineNum := tb.cursor.lineNum
	if lineNum == len(tb.lines)-1 {
		return
	}
	tb.cursor.lineNum++
	if tb.cursor.lineNum >= tb.topLine+tb.displayRange.height.GetAbsoluteValue() {
		tb.topLine++
	}
	tb.CursorOffsetAdjust()
}

func (tb *TextBuffer) GlobalOffset() int {
	byteOffset := 0
	for i := 0; i < tb.cursor.lineNum; i++ {
		length := tb.lines[i].Len()
		byteOffset += length + 1
	}
	byteOffset += tb.cursor.offset
	return byteOffset
}

func (tb *TextBuffer) UpdateLineStyle(lineNum int) {
	if lineNum > len(tb.lines)-1 {
		return
	}
	tb.lines[lineNum].UpdateLineStyle(tb.lexer)
}

func (tb *TextBuffer) Insert(r rune) {
	offset, lineNum := tb.cursor.Get()
	if offset > tb.lines[lineNum].Len() {
		return
	}
	tb.lines[lineNum].Insert(offset, byte(r))
	tb.UpdateLineStyle(lineNum)
	tb.CursorRight()
	tb.UpdateAutocompleteItems()
}

func (tb *TextBuffer) Tab() {
	insertSpaceSize := TabSize - tb.cursor.offset%TabSize
	for i := 0; i < insertSpaceSize; i++ {
		tb.Insert(' ')
	}
}

func (tb *TextBuffer) Backspace() {
	offset, lineNum := tb.cursor.Get()
	if offset > 0 {
		deletedTab := false
		tabReferenceOffset := int((offset-1)/TabSize) * TabSize
		for i := offset; i > tabReferenceOffset; i-- {
			// Notice: i-1 will always be larger than 0 because
			// tabReferenceOffset is larger or equal to 0 and i is larger than tabR.O.
			if tb.lines[lineNum].Bytes()[i-1] == ' ' {
				tb.lines[lineNum].Delete(i)
				tb.CursorLeft()
				deletedTab = true
			} else {
				break
			}
		}
		if deletedTab == false {
			tb.lines[lineNum].Delete(offset)
			tb.CursorLeft()
		}
		tb.UpdateLineStyle(lineNum)
	} else {
		if lineNum == 0 {
			return
		}
		new_offset := tb.lines[lineNum-1].Len()
		current := tb.lines[lineNum].Bytes()
		tb.lines[lineNum-1].JoinToEnd(current)
		tb.DeleteLine(lineNum)
		tb.cursor.offset = new_offset
		tb.UpdateLineStyle(lineNum - 1)
		tb.UpdateLineStyle(lineNum)
		tb.CursorUp()
	}
	tb.UpdateAutocompleteItems()
}

func (tb *TextBuffer) Return() {
	offset, lineNum := tb.cursor.Get()
	tb.InsertLine(lineNum + 1)
	extra_content := tb.lines[lineNum].CutContentToEnd(offset)
	tb.lines[lineNum+1] = NewLine()
	tb.lines[lineNum+1].Load(string(extra_content))
	tb.UpdateLineStyle(lineNum)
	tb.UpdateLineStyle(lineNum + 1)
	tb.cursor.offset = 0
	tb.CursorDown()
}

func (tb *TextBuffer) DeleteLine(lineNum int) {
	copy(tb.lines[lineNum:], tb.lines[lineNum+1:])
	tb.lines = tb.lines[:len(tb.lines)-1]
}

func (tb *TextBuffer) InsertLine(lineNum int) {
	tb.lines = append(tb.lines, nil)
	copy(tb.lines[lineNum+1:], tb.lines[lineNum:])
}

func (tb *TextBuffer) Load(url string) {
	f, err := os.Open(url)
	if err != nil {
		LogAppend(err.Error())
	}
	defer f.Close()

	tb.lexer = lexers.Match(f.Name())
	tb.fileType = tb.lexer.Config().Name
	tb.name = f.Name()
	fullDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	tb.url = fullDir + "/" + tb.name

	tb.lines = nil
	var scanner = bufio.NewScanner(f)
	for scanner.Scan() {
		content := scanner.Text()
		line := NewLine()
		line.Load(content)
		line.UpdateLineStyle(tb.lexer)
		tb.lines = append(tb.lines, line)
	}
	if tb.lines == nil {
		tb.InitFirstLine()
	}
}

func (tb *TextBuffer) Save() {
	if tb.url == "" {
		LogAppend("Cannot save unnamed file")
		return
	}
	f, err := os.Create(tb.url)
	if err != nil {
		LogAppend(err.Error())
	}
	defer f.Close()

	f.WriteString(tb.String())
}

func (tb *TextBuffer) Buffer() *bytes.Buffer {
	var b bytes.Buffer
	for _, line := range tb.lines {
		b.WriteString(line.String())
		b.WriteString("\n")
	}
	return &b
}

func (tb *TextBuffer) Bytes() []byte {
	return (*tb.Buffer()).Bytes()
}

func (tb *TextBuffer) String() string {
	return (*tb.Buffer()).String()
}

func (tb *TextBuffer) UpdateAutocompleteItems() {

	globalOffset := strconv.Itoa(tb.GlobalOffset())
	var buffer bytes.Buffer
	for _, line := range tb.lines {
		buffer.WriteString(line.String())
		buffer.WriteString("\n")
	}

	var resultBuffer bytes.Buffer

	go func() {
		c := exec.Command("gocode", "-f=json", "-source", "-builtin", "-unimported-packages", "autocomplete", tb.url, "c"+globalOffset)
		c.Stdin = &buffer
		c.Stdout = &resultBuffer
		c.Run()

		var res []interface{}
		json.Unmarshal(resultBuffer.Bytes(), &res)

		if res == nil {
			return
		}

		var completeItems []string
		for _, item := range res[1].([]interface{}) {
			completeItemName := item.(map[string]interface{})["name"].(string)
			completeItems = append(completeItems, completeItemName)
		}

		autocompleteList.LoadItems(completeItems)
		textEditor.Display()
	}()
}
