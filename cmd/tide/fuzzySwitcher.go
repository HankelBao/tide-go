package main

import (
	"os"
	"path/filepath"

	"github.com/gdamore/tcell"
	"github.com/schollz/closestmatch"
)

const (
	SwitcherTypeCommand = iota
	SwitcherTypeFiles
)

const (
	resultLineLen = 4
)

type FuzzySwitcher struct {
	displayRange *DisplayRange
	onFocus      bool

	switcherType int

	inputLine    *Line
	cursorOffset int
	cursorStart  int

	resultLines    []*Line
	selectedResult string
}

func InitFuzzySwitcher() *FuzzySwitcher {
	fuzzySwitcher := new(FuzzySwitcher)
	fuzzySwitcher.displayRange = NewDisplayRange()
	fuzzySwitcher.onFocus = false
	fuzzySwitcher.switcherType = SwitcherTypeFiles

	fuzzySwitcher.inputLine = NewLine()
	fuzzySwitcher.inputLine.Load("Files: ")
	fuzzySwitcher.inputLine.lineStyle.WithDefault(*colorTheme.boldUnderlinedStyle)
	fuzzySwitcher.cursorStart = fuzzySwitcher.inputLine.Len()
	fuzzySwitcher.cursorOffset = fuzzySwitcher.cursorStart

	fuzzySwitcher.resultLines = make([]*Line, resultLineLen)
	for i := 0; i < resultLineLen; i++ {
		fuzzySwitcher.resultLines[i] = NewLine()
	}
	return fuzzySwitcher
}

func (fs *FuzzySwitcher) Display() {

	displayLines := append([]*Line{fs.inputLine}, fs.resultLines[:]...)
	fs.displayRange.Display(displayLines, 0)
	if fs.onFocus == true {
		fs.displayRange.ShowCursor(fs.cursorOffset, 0)
	}
	screen.Show()
}

func (fs *FuzzySwitcher) Key(eventKey *tcell.EventKey) {
	key := eventKey.Key()
	switch key {
	case tcell.KeyRune:
		var inputRune = eventKey.Rune()
		fs.inputLine.Insert(fs.cursorOffset, byte(inputRune))
		fs.cursorOffset++
		fs.UpdateResults()
	case tcell.KeyUp:

	case tcell.KeyDown:

	case tcell.KeyLeft:
		if fs.cursorOffset > fs.cursorStart {
			fs.cursorOffset--
		}
	case tcell.KeyRight:
		if fs.cursorOffset < fs.inputLine.Len() {
			fs.cursorOffset++
		}
	case tcell.KeyBackspace2:
		if fs.cursorOffset > fs.cursorStart {
			fs.inputLine.Delete(fs.cursorOffset)
			fs.cursorOffset--
		}
		fs.UpdateResults()
	case tcell.KeyEnter:
		if fs.selectedResult != "" {
			textEditor.OpenFile(fs.selectedResult)
		}
	}
	fs.Display()
}

func (fs *FuzzySwitcher) Focus() {
	fs.onFocus = true
	fs.displayRange.height.AbsoluteLength(5)
	RefreshAllUIElements()
}

func (fs *FuzzySwitcher) UnFocus() {
	fs.onFocus = false
	fs.displayRange.height.AbsoluteLength(1)
	RefreshAllUIElements()
}

func (fs *FuzzySwitcher) UpdateResults() {
	fileNames, _ := ListFilesUnderDir("./")
	cm := closestmatch.New(fileNames, []int{2})
	for i, result := range cm.ClosestN(string(fs.inputLine.Bytes()[7:]), resultLineLen) {
		fs.resultLines[i].Load(result)
	}
	fs.selectedResult = fs.resultLines[0].String()
}

func ListFilesUnderDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
