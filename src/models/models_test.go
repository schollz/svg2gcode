package models

import (
	"fmt"
	"testing"
)

func TestLine(t *testing.T) {
	l := Line{}
	l.Add(1, 2)
	l.Add(-2, 3)
	l.Add(2, 4)
	ls := Lines{Lines: []Line{l}}
	ls2 := ls.BoundingBox(0, 0, 400, 400)
	ls2.Draw(-1)
	fmt.Println(ls)
	fmt.Println(ls2)
}
