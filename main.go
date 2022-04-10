package main

import (
	"fmt"
	"os"

	log "github.com/schollz/logger"
	"github.com/schollz/svg2gcode/src/gcode"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "raw",
		Version: "v0.0.1",
		Usage:   "random audio workstation",
		Commands: []*cli.Command{
			{
				Name:  "convert",
				Usage: "convert a svg to gcode",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "debug"},
					&cli.BoolFlag{Name: "animate"},
					&cli.BoolFlag{Name: "png"},
					&cli.StringFlag{Name: "in"},
					&cli.StringFlag{Name: "out"},
					&cli.Float64Flag{Name: "simplify", Usage: "in range 0-1"},
					&cli.Float64Flag{Name: "x", Required: true, Usage: "in mm"},
					&cli.Float64Flag{Name: "y", Required: true, Usage: "in mm"},
					&cli.Float64Flag{Name: "width", Required: true, Usage: "in mm"},
					&cli.Float64Flag{Name: "height", Required: true, Usage: "in mm"},
				},
				Action: func(c *cli.Context) error {
					if c.Bool("debug") {
						log.SetLevel("debug")
					}
					ops := gcode.FromSVGOptions{
						FileIn:      c.String("in"),
						FileOut:     c.String("out"),
						Animate:     c.Bool("animate"),
						PNG:         c.Bool("png"),
						BoundingBox: [4]float64{c.Float64("x"), c.Float64("y"), c.Float64("width"), c.Float64("height")},
						Simplify:    c.Float64("simplify"),
					}
					return gcode.FromSVG(ops)
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}