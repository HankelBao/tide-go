package main

type UIElement interface {
	GetContent()
	GetDisplayRange()
}

type UISelector interface {
	Key()
}