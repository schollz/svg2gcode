package svg

import (
	"fmt"
	"os"

	svg_ "github.com/rustyoz/svg"
	log "github.com/schollz/logger"
	"github.com/schollz/svg2gcode/src/models"
	"gonum.org/v1/plot/tools/bezier"
	"gonum.org/v1/plot/vg"
)

func Parse(fname string) (lines models.Lines, err error) {
	file, err := os.Open("1234.svg")
	if err != nil {
		log.Error(err)
		return
	}
	s, err := svg_.ParseSvgFromReader(file, "test", 0)
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Printf("svg: %+v\n", s.Groups[0].Elements[0])
	draw, _ := s.ParseDrawingInstructions()
	line := models.NewLine()
	var c0 vg.Point
	for {
		ins := <-draw
		if ins == nil {
			if line.Length() > 0 {
				lines.Add(line)
			}
			break
		}
		if ins.CurvePoints != nil && ins.CurvePoints.C1 != nil {
			c1 := vg.Point{X: vg.Length(ins.CurvePoints.C1[0]), Y: vg.Length(ins.CurvePoints.C1[1])}
			c2 := vg.Point{X: vg.Length(ins.CurvePoints.C2[0]), Y: vg.Length(ins.CurvePoints.C2[1])}
			c3 := vg.Point{X: vg.Length(ins.CurvePoints.T[0]), Y: vg.Length(ins.CurvePoints.T[1])}
			c := bezier.New(c0, c1, c2, c3)
			for i := 0.0; i <= 1; i += 0.1 {
				cp := c.Point(i)
				line.Add(float64(cp.X), float64(cp.Y))
			}
			c0 = c3
		} else if ins.M != nil {
			fmt.Printf("ins.M: %+v\n", ins.M)
			if line.Length() > 0 {
				lines.Add(line)
			}
			line = models.NewLine()
			line.Add(ins.M[0], ins.M[1])
			c0 = vg.Point{X: vg.Length(ins.M[0]), Y: vg.Length(ins.M[1])}
		} else {
			fmt.Printf("ins: %+v\n", ins)
		}
	}

	return
}
