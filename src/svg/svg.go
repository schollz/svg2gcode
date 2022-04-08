package svg

import (
	"fmt"
	"os"

	svg_ "github.com/rustyoz/svg"
	log "github.com/schollz/logger"
	"github.com/schollz/svg2gcode/src/models"
)

func Parse(fname string) (lines []models.Line, err error) {
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
	for {
		ins := <-draw
		if ins == nil {
			break
		}
		if ins.CurvePoints != nil && ins.CurvePoints.C1 != nil {
			line.Add(ins.CurvePoints.C1[0], ins.CurvePoints.C1[1])
			line.Add(ins.CurvePoints.C2[0], ins.CurvePoints.C2[1])
		} else if ins.M != nil {
			fmt.Printf("ins.M: %+v\n", ins.M)
			if line.Length() > 0 {
				lines = append(lines, line)
			}
			line = models.NewLine()
			line.Add(ins.M[0], ins.M[1])
		} else {
			fmt.Printf("ins: %+v\n", ins)
		}
	}
	return
}
