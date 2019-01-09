package main

import "os"

var (
	signalDoExit chan bool
	eventLoopOn  bool

	textBuffers []*TextBuffer
	colorTheme  *ColorTheme

	UIElements      []UIElement
	UISelectors     []UISelector
	focusUISelector UISelector

	statusLine *StatusLine
	textEditor *TextEditor
)

func main() {
	signalDoExit = make(chan bool)
	eventLoopOn = true

	LogInit()

	colorTheme = InitColorTheme()
	colorTheme.LoadColorTheme("emacs")

	firstTextBuffer := NewTextBuffer()
	//firstTextBuffer.Load("./test.go")

	textBuffers = append(textBuffers, firstTextBuffer)

	statusLine = InitStatusLine()
	textEditor = NewTextEditor(firstTextBuffer)

	statusLine.displayRange.width = NewRefLength([]*Length{}, WidthScale)
	statusLine.displayRange.height = NewAbsoluteLength(1)
	statusLine.displayRange.horizentalOffset = NewAbsoluteLength(0)
	statusLine.displayRange.verticalOffset = NewRefLength([]*Length{statusLine.displayRange.height}, HeightScale)
	UIElements = append(UIElements, UIElement(statusLine))

	textEditor.displayRange.width = NewRefLength([]*Length{}, WidthScale)
	textEditor.displayRange.height = NewRefLength([]*Length{statusLine.displayRange.height}, HeightScale)
	textEditor.displayRange.horizentalOffset = NewAbsoluteLength(0)
	textEditor.displayRange.verticalOffset = NewAbsoluteLength(0)
	UIElements = append(UIElements, UIElement(textEditor))
	UISelectors = append(UISelectors, UISelector(textEditor))

	focusUISelector = UISelector(textEditor)

	ScreenInit()
	go EventLoop()

	<-signalDoExit
	eventLoopOn = false
	ScreenToFini()

	LogDisplay()

	os.Exit(0)
}
