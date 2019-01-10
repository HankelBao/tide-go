package main

const LengthTypeAbsolute = true
const LengthTypeRelative = false

const WidthScale = 0
const HeightScale = 1
const NoScale = 2

type Length struct {
	absolute      bool
	absoluteValue int
	refLengths    []*Length
	refScale      int
}

func NewLength() *Length {
	length := new(Length)
	length.AbsoluteLength(0)
	return length
}

func (length *Length) AbsoluteLength(absoluteValue int) {
	length.absolute = true
	length.absoluteValue = absoluteValue
}

func (length *Length) RefLength(refLengths []*Length, refScale int) {
	length.absolute = false
	length.refLengths = refLengths
	length.refScale = refScale
}

func (l *Length) GetAbsoluteValue() int {
	if l.absolute == true {
		return l.absoluteValue
	} else {
		var extraLength int
		for _, refLength := range l.refLengths {
			if (*refLength).absolute == LengthTypeRelative {
				LogAppend("Cannot get absolute length of a relative length based on a relative length")
			}
			extraLength += (*refLength).GetAbsoluteValue()
		}

		if l.refScale == NoScale {
			return extraLength
		} else {
			var baseLength int
			if l.refScale == WidthScale {
				baseLength, _ = screen.Size()
			} else {
				_, baseLength = screen.Size()
			}
			return baseLength - extraLength
		}
	}
}
