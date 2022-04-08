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
	ls2 := ls.Normalize()
	fmt.Println(ls)
	fmt.Println(ls2)
}
