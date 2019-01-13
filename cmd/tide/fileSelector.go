package main

import (
	"os"
	"path/filepath"

	"github.com/gdamore/tcell"
)

const FileSelectorResultLen = 5

type FileSelector struct {
	displayRange *DisplayRange
	onFocus      bool

	inputLine     *InputLine
	fuzzySelector *FuzzySelector
}

var fileSelector *FileSelector

func InitFileSelector() *FileSelector {
	fileSelector := new(FileSelector)
	fileSelector.displayRange = NewDisplayRange()
	fileSelector.onFocus = false
	fileSelector.inputLine = NewInputLine("Files> ")
	fileSelector.fuzzySelector = NewFuzzySelector(FileSelectorResultLen, UIElement(fileSelector))
	go func() {
		filesRange, _ := ListFilesUnderDir("./")
		fileSelector.fuzzySelector.UpdateRange(filesRange)
	}()
	return fileSelector
}

func (fs *FileSelector) Display() {
	if fs.onFocus == false {
		return
	}
	displayLines := append([]*Line{fs.inputLine.GetLine()}, fs.fuzzySelector.resultLines...)
	fs.displayRange.Display(displayLines, 0)
	fs.displayRange.ShowCursor(fs.inputLine.cursorOffset, 0)
	screen.Show()
}

func (fs *FileSelector) Key(eventKey *tcell.EventKey) {
	key := eventKey.Key()
	switch key {
	case tcell.KeyRune:
		var inputRune = eventKey.Rune()
		fs.inputLine.Insert(inputRune)
		fs.Display()
		fs.fuzzySelector.UpdateResults(fs.inputLine.GetInput())
	case tcell.KeyLeft:
		fs.inputLine.CursorLeft()
		fs.Display()
	case tcell.KeyRight:
		fs.inputLine.CursorRight()
		fs.Display()
	case tcell.KeyBackspace2:
		fs.inputLine.Backspace()
		fs.Display()
		fs.fuzzySelector.UpdateResults(fs.inputLine.GetInput())
	case tcell.KeyEnter:
		fileName := fs.fuzzySelector.GetTopResult()
		if fileName != "" {
			textEditor.OpenFile(fileName)
			SwitchUIFocus(textEditor)
		}
	}
}

func (fs *FileSelector) Focus() {
	fs.onFocus = true
	fs.displayRange.height.AbsoluteLength(FileSelectorResultLen + 1)
	RefreshAllUIElements()
}

func (fs *FileSelector) UnFocus() {
	fs.onFocus = false
	fs.displayRange.height.AbsoluteLength(0)
	RefreshAllUIElements()
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
