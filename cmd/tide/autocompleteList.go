package main

const MaxAutocompleteListDisplayLen = 6
const AutocompleteListDisplayWidth = 30

type AutocompleteList struct {
	displayRange *DisplayRange
	visible      bool

	completeOffset int
	completeItems  []CompleteItem
	displayLines   []*Line
}

func InitAutocompleteList() *AutocompleteList {
	ac := new(AutocompleteList)
	ac.displayRange = NewDisplayRange()
	ac.displayRange.horizentalOffset.AbsoluteLength(0)
	ac.displayRange.verticalOffset.AbsoluteLength(0)
	ac.displayRange.height.AbsoluteLength(0)
	ac.displayRange.width.AbsoluteLength(AutocompleteListDisplayWidth)
	ac.visible = false
	ac.displayLines = make([]*Line, MaxAutocompleteListDisplayLen)
	for i := 0; i < MaxAutocompleteListDisplayLen; i++ {
		ac.displayLines[i] = NewLine()
		ac.displayLines[i].lineStyle.WithDefault(*colorTheme.boldStyle)
	}
	ac.completeItems = []CompleteItem{}
	return ac
}

func (ac *AutocompleteList) LoadItems(completeOffset int, items []CompleteItem) {
	ac.completeOffset = completeOffset
	ac.completeItems = items
	if len(items) == 0 {
		ac.visible = false
	} else {
		ac.visible = true
	}
}

func (ac *AutocompleteList) Display() {
	if ac.visible == false {
		return
	}

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
			ac.displayLines[i].Load(ac.completeItems[i].name)
		}
	}

	xPos, yPos := textEditor.displayRange.GetScreenPos(textEditor.textBuffer.cursor.Get())
	ac.displayRange.horizentalOffset.AbsoluteLength(xPos)
	ac.displayRange.verticalOffset.AbsoluteLength(yPos + 1)
	ac.displayRange.height.AbsoluteLength(autocompleteListLen)
	ac.displayRange.Display(ac.displayLines[:autocompleteListLen], 0)
	screen.Show()
}

func (ac *AutocompleteList) GetComplete() string {
	completeItem := ac.completeItems[0]
	hintLine.LoadContent(completeItem.description)
	return completeItem.name[ac.completeOffset:]
}
