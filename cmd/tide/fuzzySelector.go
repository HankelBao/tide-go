package main

import "github.com/schollz/closestmatch"

type FuzzySelector struct {
	rangeString []string

	// The result lines here should be used directly by user
	resultLines     []*Line
	targetUIElement UIElement
}

func NewFuzzySelector(resultLineLen int, targetUIElement UIElement) *FuzzySelector {
	fs := new(FuzzySelector)
	fs.resultLines = make([]*Line, resultLineLen)
	for i := 0; i < resultLineLen; i++ {
		fs.resultLines[i] = NewLine()
	}

	fs.targetUIElement = targetUIElement
	return fs
}

func (fs *FuzzySelector) UpdateRange(rangeString []string) {
	fs.rangeString = rangeString
}

func (fs *FuzzySelector) UpdateResults(input string) {
	go func() {
		cm := closestmatch.New(fs.rangeString, []int{2})
		resultLen := len(fs.resultLines)
		for i, result := range cm.ClosestN(input, resultLen) {
			fs.resultLines[i].Load(result)
		}
		fs.targetUIElement.Display()
	}()
}

func (fs *FuzzySelector) GetTopResult() string {
	return fs.resultLines[0].String()
}
