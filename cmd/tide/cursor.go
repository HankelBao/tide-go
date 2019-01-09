package main

type Cursor struct {
	lineNum int
	offset  int
}

func (c *Cursor) Get() (int, int) {
	return c.lineNum, c.offset
}
