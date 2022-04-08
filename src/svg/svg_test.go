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
	b, _ := json.MarshalIndent(lines, "", " ")
	fmt.Println(string(b))
}
