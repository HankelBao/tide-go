package main

import "os"

var (
	signalDoExit chan bool
	eventLoopOn  bool

	textBuffers []*TextBuffer
	colorTheme  *ColorTheme

	UISelectors     []UISelector
	focusUISelector UISelector

	statusLine    *StatusLine
	textEditor    *TextEditor
	fuzzySwitcher *FuzzySwitcher
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
	fuzzySwitcher = InitFuzzySwitcher()
	textEditor = InitTextEditor(firstTextBuffer)

	fuzzySwitcher.displayRange.width.RefLength([]*Length{}, WidthScale)
	fuzzySwitcher.displayRange.height.AbsoluteLength(1)
	fuzzySwitcher.displayRange.horizentalOffset.AbsoluteLength(0)
	fuzzySwitcher.displayRange.verticalOffset.RefLength([]*Length{fuzzySwitcher.displayRange.height}, HeightScale)

	statusLine.displayRange.width.RefLength([]*Length{}, WidthScale)
	statusLine.displayRange.height.AbsoluteLength(1)
	statusLine.displayRange.horizentalOffset.AbsoluteLength(0)
	statusLine.displayRange.verticalOffset.RefLength([]*Length{fuzzySwitcher.displayRange.height, statusLine.displayRange.height}, HeightScale)

	textEditor.displayRange.width.RefLength([]*Length{}, WidthScale)
	textEditor.displayRange.height.RefLength([]*Length{statusLine.displayRange.height, fuzzySwitcher.displayRange.height}, HeightScale)
	textEditor.displayRange.horizentalOffset.AbsoluteLength(0)
	textEditor.displayRange.verticalOffset.AbsoluteLength(0)
	UISelectors = append(UISelectors, UISelector(textEditor))

	ScreenInit()

	RefreshAllUIElements()

	focusUISelector = UISelector(textEditor)
	focusUISelector.Focus()

	go EventLoop()

	<-signalDoExit
	eventLoopOn = false
	ScreenToFini()

	LogDisplay()

	os.Exit(0)
}
