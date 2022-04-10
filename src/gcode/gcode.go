package gcode

import (
	"io/ioutil"
	"strings"

	log "github.com/schollz/logger"
	"github.com/schollz/svg2gcode/src/svg"
	"github.com/tarm/serial"
	"go.bug.st/serial/enumerator"
)

type FromSVGOptions struct {
	FileIn      string
	FileOut     string
	BoundingBox [4]float64
	Simplify    float64
	Animate     bool
	PNG         bool
	MaxLenth    float64
	Consolidate float64
}

func FromSVG(op FromSVGOptions) (err error) {
	if op.FileOut == "" {
		op.FileOut = op.FileIn
	}
	lines, err := svg.ParseSVG(op.FileIn)
	if err != nil {
		return
	}
	lines = lines.BoundingBox(op.BoundingBox[0], op.BoundingBox[1], op.BoundingBox[2], op.BoundingBox[3])
	lines = lines.Consolidate(op.Consolidate)
	lines = lines.RemoveSmall(op.MaxLenth)
	lines = lines.BestOrdering()
	lines = lines.Consolidate(op.Consolidate)
	lines = lines.Simplify(op.Simplify)
	if op.Animate {
		log.Debugf("animating %s", op.FileOut+".gif")
		lines.Animate(op.FileOut + ".gif")
	}
	log.Debugf("have %d lines", len(lines.Lines))
	if op.PNG {
		log.Debugf("drawing %s", op.FileOut+".png")
		lines.Draw(op.FileOut + ".png")
	}
	err = ioutil.WriteFile(op.FileOut+".gcode", []byte(lines.ToGcode()), 0644)
	return
}

func Send(instructions string, portName0 ...string) (err error) {
	portName := ""
	if len(portName0) > 0 {
		portName = portName0[0]
	} else {
		var ports []*enumerator.PortDetails
		ports, err = enumerator.GetDetailedPortsList()
		if err != nil {
			log.Error(err)
			return
		}
		if len(ports) == 0 {
			return
		}
		for _, port := range ports {
			if strings.Contains(port.Product, "CH340") {
				portName = port.Name
			}
		}
	}

	if portName == "" {
		log.Error("could not find port for pritner")
		return
	}
	c := &serial.Config{Name: portName, Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Error(err)
		return
	}
	defer s.Close()

	for _, instruction := range strings.Split(instructions, "\n") {
		instruction = strings.TrimSpace(instruction)
		if instruction == "" {
			continue
		}
		log.Debug(instruction)
		_, err = s.Write([]byte(instruction + "\n"))
		if err != nil {
			log.Errorf("%s: %s", instruction, err.Error())
			return
		}
	}
	return
}
