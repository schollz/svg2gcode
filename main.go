package main

import (
	"fmt"
	"io/ioutil"
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
					&cli.Float64Flag{Name: "simplify", Usage: "as %% (0-1)"},
					&cli.Float64Flag{Name: "min-length", Usage: "as %% of size (0-1)"},
					&cli.Float64Flag{Name: "consolidate", Usage: "consolidates lines %% of size (0-1)"},
					&cli.Float64Flag{Name: "x", Usage: "in mm"},
					&cli.Float64Flag{Name: "y", Usage: "in mm"},
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
						MinLength:   c.Float64("min-length"),
						Consolidate: c.Float64("consolidate"),
					}
					return gcode.FromSVG(ops)
				},
			}, {
				Name:  "upload",
				Usage: "upload gcode to a cnc machine or 3d printer",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "debug"},
					&cli.StringFlag{Name: "in"},
				},
				Action: func(c *cli.Context) error {
					if c.Bool("debug") {
						log.SetLevel("debug")
					}
					instructions, err := ioutil.ReadFile(c.String("in"))
					if err != nil {
						return err
					}
					return gcode.Send(string(instructions))
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
