package svg

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	lines, err := Parse("1234.svg")
	if err != nil {
		t.Error(err)
	}
	ls2 := lines.BoundingBox(0, 0, 100, 100)
	ls2.Animate("test.gif")
	ls2.Draw("final.png")

	b, _ := json.MarshalIndent(lines, "", " ")
	fmt.Println(string(b))
}
