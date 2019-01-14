package main

type DisplayRange struct {
	horizentalOffset *Length
	verticalOffset   *Length
	width            *Length
	height           *Length
}

func NewDisplayRange() *DisplayRange {
	displayRange := new(DisplayRange)
	displayRange.horizentalOffset = NewLength()
	displayRange.verticalOffset = NewLength()
	displayRange.width = NewLength()
	displayRange.height = NewLength()
	return displayRange
}

func (dr *DisplayRange) Display(lines []*Line, lineOffset int) {
	xOffset := dr.horizentalOffset.GetAbsoluteValue()
	yOffset := dr.verticalOffset.GetAbsoluteValue()
	width := dr.width.GetAbsoluteValue()
	height := dr.height.GetAbsoluteValue()

	for i := 0; i < height; i++ {
		for j := lineOffset; j <= lineOffset+width; j++ {
			realX, realY := xOffset+j-lineOffset, yOffset+i
			if i < len(lines) {
				if j < lines[i].Len() {
					r := rune(lines[i].Bytes()[j])
					style := lines[i].lineStyle.GetStyleAt(j)
					screen.SetContent(realX, realY, r, nil, style)
					continue
				}
				screen.SetContent(realX, realY, ' ', nil, lines[i].lineStyle.GetDefaultStyle())
				continue
			}
			screen.SetContent(realX, realY, ' ', nil, *colorTheme.defaultStyle)
		}
	}
}

func (dr *DisplayRange) GetScreenPos(xOffset int, yOffset int) (int, int) {
	return dr.horizentalOffset.GetAbsoluteValue() + xOffset, dr.verticalOffset.GetAbsoluteValue() + yOffset
}

func (dr *DisplayRange) ShowCursor(x int, y int) {
	realX := dr.horizentalOffset.GetAbsoluteValue() + x
	realY := dr.verticalOffset.GetAbsoluteValue() + y
	screen.ShowCursor(realX, realY)
}
