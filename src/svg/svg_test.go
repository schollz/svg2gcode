package svg

import (
	"testing"
)

// func TestLine(t *testing.T) {
// 	l := Line{}
// 	l.Add(1, 2)
// 	l.Add(-2, 3)
// 	l.Add(2, 4)
// 	ls := Lines{Lines: []Line{l}}
// 	ls2 := ls.BoundingBox(0, 0, 400, 400)
// 	ls2.Draw(-1)
// 	fmt.Println(ls)
// 	fmt.Println(ls2)
// }

func TestParse(t *testing.T) {
	lines, err := ParseSVG("1234.svg")
	if err != nil {
		t.Error(err)
	}
	ls2 := lines.BoundingBox(0, 0, 300, 300)
	ls2 = ls2.Simplify(0.55)
	ls2.Animate("test.gif")
	ls2.Draw("final.png")
}
