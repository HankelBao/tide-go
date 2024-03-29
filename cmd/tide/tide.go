package main

import "os"

// Global Variables
var (
	signalDoExit chan bool
	eventLoopOn  bool

	textBuffers []*TextBuffer
	colorTheme  *ColorTheme

	focusUISelector UISelector

	statusLine       *StatusLine
	textEditor       *TextEditor
	autocompleteList *AutocompleteList
	hintLine         *HintLine
	fileSelector     *FileSelector
)

func main() {
	signalDoExit = make(chan bool)
	eventLoopOn = true

	LogInit()

	colorTheme = InitColorTheme()
	colorTheme.LoadColorTheme("emacs")

	firstTextBuffer := NewTextBuffer()

	textBuffers = append(textBuffers, firstTextBuffer)

	statusLine = InitStatusLine()
	textEditor = InitTextEditor(firstTextBuffer)
	autocompleteList = InitAutocompleteList()
	hintLine = InitHintLine()
	fileSelector = InitFileSelector()

	fileSelector.displayRange.width.RefLength([]*Length{}, WidthScale)
	fileSelector.displayRange.height.AbsoluteLength(0)
	fileSelector.displayRange.horizentalOffset.AbsoluteLength(0)
	fileSelector.displayRange.verticalOffset.RefLength([]*Length{fileSelector.displayRange.height}, HeightScale)

	statusLine.displayRange.width.RefLength([]*Length{}, WidthScale)
	statusLine.displayRange.height.AbsoluteLength(1)
	statusLine.displayRange.horizentalOffset.AbsoluteLength(0)
	statusLine.displayRange.verticalOffset.RefLength([]*Length{fileSelector.displayRange.height, statusLine.displayRange.height}, HeightScale)

	textEditor.displayRange.width.RefLength([]*Length{}, WidthScale)
	textEditor.displayRange.height.RefLength([]*Length{statusLine.displayRange.height, fileSelector.displayRange.height}, HeightScale)
	textEditor.displayRange.horizentalOffset.AbsoluteLength(0)
	textEditor.displayRange.verticalOffset.AbsoluteLength(0)

	ScreenInit()

	RefreshAllUIElements()

	InitUIFocus(textEditor)

	go EventLoop()

	<-signalDoExit
	eventLoopOn = false
	ScreenToFini()

	LogDisplay()

	os.Exit(0)
}
