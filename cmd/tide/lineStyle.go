package main

import (
	"github.com/alecthomas/chroma"
	"github.com/gdamore/tcell"
)

type StyleDescriptor struct {
	size  int
	style *tcell.Style
}

type LineStyle struct {
	styleDescriptors []*StyleDescriptor
	defaultStyle     tcell.Style
}

func NewLineStyle() *LineStyle {
	lineStyle := new(LineStyle)
	lineStyle.WithDefault(*colorTheme.defaultStyle)
	return lineStyle
}

func (lineStyle *LineStyle) LoadFromSourceCode(sourceCode string, lexer chroma.Lexer) {
	iterator, err := lexer.Tokenise(nil, sourceCode)
	if err != nil {
		LogAppend("Convert Line to Token failed")
	}
	lineStyle.styleDescriptors = make([]*StyleDescriptor, 0)
	for _, token := range iterator.Tokens() {
		styleDescriptor := new(StyleDescriptor)
		styleDescriptor.size = len(token.Value)
		styleDescriptor.style = colorTheme.GetStyle(token.Type)
		lineStyle.styleDescriptors = append(lineStyle.styleDescriptors, styleDescriptor)
	}
}

func (lineStyle *LineStyle) WithDefault(defaultStyle tcell.Style) {
	lineStyle.styleDescriptors = nil
	lineStyle.defaultStyle = defaultStyle
}

func (ls *LineStyle) GetStyleAt(offset int) tcell.Style {
	realOffset := 0
	for _, styleDescriptor := range ls.styleDescriptors {
		realOffset += styleDescriptor.size
		if realOffset > offset {
			return *styleDescriptor.style
		}
	}
	return ls.defaultStyle
}

func (ls *LineStyle) GetDefaultStyle() tcell.Style {
	return ls.defaultStyle
}
