package main

const MaxAutocompleteListDisplayLen = 6
const AutocompleteListDisplayWidth = 20

type AutocompleteList struct {
	displayRange *DisplayRange
	onFocus      bool

	completeItems []string
	displayLines  []*Line
}

func InitAutocompleteList() *AutocompleteList {
	ac := new(AutocompleteList)
	ac.displayRange = NewDisplayRange()
	ac.displayRange.horizentalOffset.AbsoluteLength(0)
	ac.displayRange.verticalOffset.AbsoluteLength(0)
	ac.displayRange.height.AbsoluteLength(0)
	ac.displayRange.width.AbsoluteLength(AutocompleteListDisplayWidth)
	ac.onFocus = false
	ac.displayLines = make([]*Line, MaxAutocompleteListDisplayLen)
	for i := 0; i < MaxAutocompleteListDisplayLen; i++ {
		ac.displayLines[i] = NewLine()
		ac.displayLines[i].lineStyle.WithDefault(*colorTheme.boldStyle)
	}
	ac.completeItems = []string{}
	return ac
}

func (ac *AutocompleteList) LoadItems(items []string) {
	ac.completeItems = items
}

func (ac *AutocompleteList) Display() {
	var autocompleteListLen int
	if len(ac.completeItems) > MaxAutocompleteListDisplayLen {
		autocompleteListLen = MaxAutocompleteListDisplayLen
	} else {
		autocompleteListLen = len(ac.completeItems)
	}

	for i := 0; i < autocompleteListLen; i++ {
		if i >= len(ac.completeItems) {
			ac.displayLines[i].Load("")
		} else {
			ac.displayLines[i].Load(ac.completeItems[i])
		}
	}

	xPos, yPos := textEditor.displayRange.GetScreenPos(textEditor.textBuffer.cursor.Get())
	ac.displayRange.horizentalOffset.AbsoluteLength(xPos)
	ac.displayRange.verticalOffset.AbsoluteLength(yPos)
	ac.displayRange.height.AbsoluteLength(autocompleteListLen)
	ac.displayRange.Display(ac.displayLines[:autocompleteListLen], 0)
	screen.Show()
}
