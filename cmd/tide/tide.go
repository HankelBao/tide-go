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

	bufferline *BufferLine
	textEditor *TextEditor
)

func main() {
	signalDoExit = make(chan bool)
	eventLoopOn = true

	LogInit()

	colorTheme = InitColorTheme()
	colorTheme.LoadColorTheme("emacs")

	firstTextBuffer := NewTextBuffer()
	firstTextBuffer.Load("./test.go")

	textBuffers = append(textBuffers, firstTextBuffer)

	bufferline = NewBufferLine()
	textEditor = NewTextEditor(firstTextBuffer)

	bufferline.displayRange.horizentalOffset = NewAbsoluteLength(0)
	bufferline.displayRange.verticalOffset = NewAbsoluteLength(0)
	bufferline.displayRange.width = NewRefLength([]*Length{}, WidthScale)
	bufferline.displayRange.height = NewAbsoluteLength(2)
	UIElements = append(UIElements, UIElement(bufferline))

	textEditor.displayRange.horizentalOffset = NewAbsoluteLength(0)
	textEditor.displayRange.verticalOffset = NewRefLength([]*Length{bufferline.displayRange.height}, NoScale)
	textEditor.displayRange.width = NewRefLength([]*Length{}, WidthScale)
	textEditor.displayRange.height = NewRefLength([]*Length{bufferline.displayRange.height}, HeightScale)
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
