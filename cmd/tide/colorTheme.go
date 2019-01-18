package main

import (
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/styles"
	"github.com/gdamore/tcell"
)

type ColorTheme struct {
	styleMap map[chroma.TokenType]*tcell.Style

	defaultStyle        *tcell.Style
	reversedStyle       *tcell.Style
	underlinedStyle     *tcell.Style
	boldStyle           *tcell.Style
	boldUnderlinedStyle *tcell.Style
}

func InitColorTheme() *ColorTheme {
	colorTheme := new(ColorTheme)
	colorTheme.styleMap = make(map[chroma.TokenType]*tcell.Style)
	return colorTheme
}

func (ct *ColorTheme) LoadColorTheme(themeName string) bool {
	styleManager := styles.Get(themeName)
	if styleManager == nil {
		LogAppend("Load unknown color theme")
		return false
	}
	for tt, _ := range chroma.StandardTypes {
		tcellStyle := ConvertStyleEntryToTCellStyle(styleManager.Get(tt))
		ct.styleMap[tt] = tcellStyle
	}

	ct.defaultStyle = ConvertStyleEntryToTCellStyle(styleManager.Get(chroma.Text))
	reversedStyle := (*ct.defaultStyle).Reverse(true)
	underlinedStyle := (*ct.defaultStyle).Underline(true)
	boldStyle := (*ct.defaultStyle).Bold(true)
	boldUnderlinedStyle := boldStyle.Underline(true)

	ct.reversedStyle = &reversedStyle
	ct.underlinedStyle = &underlinedStyle
	ct.boldStyle = &boldStyle
	ct.boldUnderlinedStyle = &boldUnderlinedStyle

	return true
}

func ConvertStyleEntryToTCellStyle(se chroma.StyleEntry) *tcell.Style {
	var tcellStyle tcell.Style
	tcellStyle = tcellStyle.Background(tcell.GetColor(se.Background.String()))
	tcellStyle = tcellStyle.Foreground(tcell.GetColor(se.Colour.String()))
	return &tcellStyle
}

func (ct *ColorTheme) GetStyle(tokenType chroma.TokenType) *tcell.Style {
	return ct.styleMap[tokenType]
}
