package svg

import (
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"math"
	"os"
	"runtime"
	"sort"

	"github.com/fogleman/gg"
	svg_ "github.com/rustyoz/svg"
	log "github.com/schollz/logger"
	"github.com/schollz/progressbar/v3"
	"github.com/schollz/svg2gcode/src/grad"
	"gonum.org/v1/plot/tools/bezier"
	"gonum.org/v1/plot/vg"
)

const MillimetersPerPixel = 0.26458333

type Lines struct {
	Lines  []Line
	Height float64
	Width  float64
	Bounds [4]float64
}

type Line struct {
	Points []Point
}

type Point struct {
	X, Y float64
}

func Distance(p1, p2 Point) float64 {
	return math.Sqrt(math.Pow(p1.X-p2.X, 2) + math.Pow(p1.Y-p2.Y, 2))
}

func ParseSVG(fname string) (lines Lines, err error) {
	file, err := os.Open(fname)
	if err != nil {
		log.Error(err)
		return
	}
	s, err := svg_.ParseSvgFromReader(file, "test", 0)
	if err != nil {
		log.Error(err)
		return
	}
	draw, _ := s.ParseDrawingInstructions()
	line := NewLine()
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
			if line.Length() > 0 {
				lines.Add(line)
			}
			line = NewLine()
			line.Add(ins.M[0], ins.M[1])
			c0 = vg.Point{X: vg.Length(ins.M[0]), Y: vg.Length(ins.M[1])}
		} else {
		}
	}

	return
}

func (l Lines) Copy() (l2 Lines) {
	l2 = Lines{Bounds: l.Bounds, Height: l.Height, Width: l.Width}
	for _, line := range l.Lines {
		line2 := Line{}
		for _, p := range line.Points {
			line2.Points = append(line2.Points, Point{p.X, p.Y})
		}
		l2.Lines = append(l2.Lines, line2)
	}
	return
}

func (l Lines) ToGcode() (gcode string) {
	for _, line := range l.Lines {
		gcode += "G91 ; Set coordinates to relative\n"
		gcode += "G1 Z10 F1000 ; raise pen\n"
		gcode += "G90 ; Set coordinates to absolute\n"
		for j, p := range line.Points {
			gcode += fmt.Sprintf("G1 X%2.3f Y%2.3f F1000\n", p.X, p.Y)
			if j == 0 {
				gcode += "G91 ; Set coordinates to relative\n"
				gcode += "G1 Z-10 F1000 ; lower pen\n"
				gcode += "G90 ; Set coordinates to absolute\n"
			}
		}
	}
	return
}

// Simplify points by some percentage using the
// https://en.wikipedia.org/wiki/Visvalingam%E2%80%93Whyatt_algorithm
func (l Lines) Simplify(percent float64) (l2 Lines) {
	if percent == 0 {
		l2 = l.Copy()
		return
	}
	pointsAvailable := 0.0
	for _, line := range l.Lines {
		if len(line.Points) > 2 {
			pointsAvailable += float64(len(line.Points) - 2)
		}
	}

	type Area struct {
		val    float64
		linei  int
		pointi int
	}
	areas := []Area{}
	for linei, line := range l.Lines {
		if len(line.Points) <= 2 {
			continue
		}
		for i := 1; i < len(line.Points)-1; i++ {
			x0 := line.Points[i-1].X
			y0 := line.Points[i-1].Y
			x1 := line.Points[i].X
			y1 := line.Points[i].Y
			x2 := line.Points[i+1].X
			y2 := line.Points[i+1].Y
			area := Area{
				val:    0.5 * (x0*y1 + x1*y2 + x2*y0 - x0*y2 - x1*y0 - x2*y1),
				linei:  linei,
				pointi: i,
			}
			areas = append(areas, area)
		}
	}
	sort.Slice(areas, func(i, j int) bool {
		return areas[i].val < areas[j].val
	})

	removal := make(map[string]struct{})
	for i := 0; i < int(percent*pointsAvailable); i++ {
		removal[fmt.Sprintf("%d-%d", areas[i].linei, areas[i].pointi)] = struct{}{}
	}

	l2 = Lines{Bounds: l.Bounds, Height: l.Height, Width: l.Width}
	for linei, line := range l.Lines {
		line2 := Line{}
		for pointi, p := range line.Points {
			if _, ok := removal[fmt.Sprintf("%d-%d", linei, pointi)]; !ok {
				line2.Points = append(line2.Points, Point{p.X, p.Y})
			}
		}
		l2.Lines = append(l2.Lines, line2)
	}

	pointsNow := 0.0
	for _, line := range l2.Lines {
		if len(line.Points) > 2 {
			pointsNow += float64(len(line.Points) - 2)
		}
	}

	log.Debugf("simplified %2.3f%%: %2.0f points -> %2.0f points", percent*100, pointsAvailable, pointsNow)

	return
}

func (l *Lines) Add(line Line) {
	l.Lines = append(l.Lines, line)
}

func (l Lines) Draw(fname string) (err error) {
	im := l.DrawStep(-1)
	f, err := os.Create(fname)
	if err != nil {
		log.Error(err)
		return
	}
	defer f.Close()

	err = png.Encode(f, im)
	if err != nil {
		log.Error(err)
		return
	}
	return
}

func (l Lines) DrawStep(step int, im0 ...image.Image) (im image.Image) {
	var dc *gg.Context
	if len(im0) > 0 {
		dc = gg.NewContextForImage(im0[0])
	} else {
		dc = gg.NewContext(int((l.Bounds[0]+l.Bounds[2])/MillimetersPerPixel), int((l.Bounds[1]+l.Bounds[3])/MillimetersPerPixel))
		dc.DrawRectangle(0, 0, (l.Bounds[0]+l.Bounds[2])/MillimetersPerPixel, (l.Bounds[1]+l.Bounds[3])/MillimetersPerPixel)
		dc.SetRGB(1, 1, 1)
		dc.Fill()
	}
	i := 0
	for _, line := range l.Lines {
		for j, point := range line.Points {
			if j == 0 {
				continue
			}
			dc.DrawLine(line.Points[j-1].X/MillimetersPerPixel, line.Points[j-1].Y/MillimetersPerPixel, point.X/MillimetersPerPixel, point.Y/MillimetersPerPixel)
			if i == step && step > -1 {
				break
			}
			i++
		}
		if i == step && step > -1 {
			break
		}
	}
	dc.SetRGB(255, 255, 255)
	dc.SetLineWidth(2)
	dc.Stroke()
	im = dc.Image()
	return
}

func (l Lines) Animate(fname string) (err error) {
	numJobs := 0
	for _, line := range l.Lines {
		for j := range line.Points {
			if j == 0 {
				continue
			}
			numJobs++
		}
	}

	type job struct {
		simage image.Image
		i      int
	}
	type result struct {
		palettedImage *image.Paletted
		i             int
	}

	jobs := make(chan job, numJobs)
	results := make(chan result, numJobs)
	runtime.GOMAXPROCS(runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		go func(jobs <-chan job, results chan<- result) {
			for j := range jobs {
				// step 3: specify the work for the worker
				var r result
				bounds := j.simage.Bounds()
				r.i = j.i
				r.palettedImage = image.NewPaletted(bounds, palette.WebSafe)
				draw.Draw(r.palettedImage, r.palettedImage.Rect, j.simage, bounds.Min, draw.Over)
				results <- r
			}
		}(jobs, results)
	}

	// step 4: send out jobs
	var simage image.Image
	for i := 0; i < numJobs; i++ {
		if i == 0 {
			simage = l.DrawStep(0)
		} else {
			simage = l.DrawStep(i, simage)
		}
		jobs <- job{simage, i}
	}
	close(jobs)

	// step 5: do something with results
	palettedImages := make([]*image.Paletted, numJobs)
	bar := progressbar.Default(int64(numJobs))
	for i := 0; i < numJobs; i++ {
		bar.Add(1)
		r := <-results
		palettedImages[r.i] = r.palettedImage
	}

	outGif := &gif.GIF{}
	for i := 0; i < numJobs; i++ {
		outGif.Image = append(outGif.Image, palettedImages[i])
		if i < numJobs-1 {
			outGif.Delay = append(outGif.Delay, 0)
		} else {
			outGif.Delay = append(outGif.Delay, 200)
		}
	}

	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return
	}
	err = gif.EncodeAll(f, outGif)
	if err != nil {
		return
	}
	f.Close()
	return
}

func (l Lines) BoundingBox(x, y, w, h float64) (l2 Lines) {
	l2 = l.Normalize()
	multiply := w
	if l2.Height > l2.Width {
		multiply = h
	}

	maxPoint := Point{-1000000, -1000000}
	for i, line := range l2.Lines {
		for j := range line.Points {
			l2.Lines[i].Points[j].X *= multiply
			if l2.Lines[i].Points[j].X > maxPoint.X {
				maxPoint.X = l2.Lines[i].Points[j].X
			}
			l2.Lines[i].Points[j].Y *= multiply
			if l2.Lines[i].Points[j].Y > maxPoint.Y {
				maxPoint.Y = l2.Lines[i].Points[j].Y
			}
			l2.Lines[i].Points[j].X += x
			l2.Lines[i].Points[j].Y += y
		}
	}
	l2.Bounds = [4]float64{x, y, maxPoint.X, maxPoint.Y}

	return
}

// Normalize keeps aspect ratio but makes sure all coordinates are between 0 and 1
func (l Lines) Normalize() (l2 Lines) {
	l2 = Lines{}
	for _, line := range l.Lines {
		line2 := Line{}
		for _, p := range line.Points {
			line2.Points = append(line2.Points, Point{p.X, p.Y})
		}
		l2.Lines = append(l2.Lines, line2)
	}
	// first make sure everything is above 0
	minPoint := Point{1000000000, 1000000000}
	maxPoint := Point{-1000000, -1000000}
	for _, line := range l2.Lines {
		for _, p := range line.Points {
			if p.X < minPoint.X {
				minPoint.X = p.X
			}
			if p.X > maxPoint.X {
				maxPoint.X = p.X
			}
			if p.Y < minPoint.Y {
				minPoint.Y = p.Y
			}
			if p.Y > maxPoint.Y {
				maxPoint.Y = p.Y
			}
		}
	}
	_ = maxPoint
	for i, line := range l2.Lines {
		for j := range line.Points {
			l2.Lines[i].Points[j].X -= minPoint.X
			l2.Lines[i].Points[j].Y -= minPoint.Y
		}
	}

	// find the biggest part and reduce everything by that
	biggestNum := 0.0
	for _, line := range l2.Lines {
		for _, p := range line.Points {
			if p.X > biggestNum {
				biggestNum = p.X
			}
			if p.X > l2.Width {
				l2.Width = p.X
			}
			if p.Y > biggestNum {
				biggestNum = p.Y
			}
			if p.Y > l2.Height {
				l2.Height = p.Y
			}
		}
	}
	for i, line := range l2.Lines {
		for j := range line.Points {
			l2.Lines[i].Points[j].X /= biggestNum
			l2.Lines[i].Points[j].Y /= biggestNum
		}
	}
	l2.Height /= biggestNum
	l2.Width /= biggestNum
	return
}

func (l Line) PathLength() (pathLength float64) {
	for i := 1; i < len(l.Points); i++ {
		pathLength += Distance(l.Points[i-1], l.Points[i])
	}
	return
}

func (l Lines) RemoveSmall(fraction float64) (l2 Lines) {
	dist := fraction * (l.Bounds[2] + l.Bounds[3]) / 2
	l2 = Lines{Bounds: l.Bounds, Height: l.Height, Width: l.Width}
	for _, line := range l.Lines {
		if line.PathLength() < dist {
			continue
		}
		line2 := Line{}
		for _, p := range line.Points {
			line2.Points = append(line2.Points, Point{p.X, p.Y})
		}
		l2.Lines = append(l2.Lines, line2)
	}
	log.Debugf("cleaned %d lines to %d lines", len(l.Lines), len(l2.Lines))
	return
}

func (l Lines) BestOrdering() (l2 Lines) {
	points := []grad.Point{}
	for _, line := range l.Lines {
		points = append(points, grad.Point{line.Points[0].X, line.Points[0].Y})
		points = append(points, grad.Point{line.Points[len(line.Points)-1].X, line.Points[len(line.Points)-1].Y})
	}
	path := grad.FindBestPath(points)
	l2 = Lines{Bounds: l.Bounds, Height: l.Height, Width: l.Width}
	for i := 0; i < len(path); i += 2 {
		linei := int(math.Floor(float64(path[i]) / 2.0))
		line := l.Lines[linei]
		line2 := Line{}
		if path[i] > path[i+1] {
			// reverse points
			for pi := len(line.Points) - 1; pi >= 0; pi-- {
				p := line.Points[pi]
				line2.Points = append(line2.Points, Point{p.X, p.Y})
			}
		} else {
			for _, p := range line.Points {
				line2.Points = append(line2.Points, Point{p.X, p.Y})
			}
		}
		l2.Lines = append(l2.Lines, line2)
	}
	return
}

func (l Lines) Consolidate(joinFraction float64) (l2 Lines) {
	l2 = Lines{Bounds: l.Bounds, Height: l.Height, Width: l.Width}
	for _, line := range l.Lines {
		line2 := Line{}
		for _, p := range line.Points {
			line2.Points = append(line2.Points, Point{p.X, p.Y})
		}
		if len(l2.Lines) > 0 {
			lastPoint := l2.Lines[len(l2.Lines)-1].Points[len(l2.Lines[len(l2.Lines)-1].Points)-1]
			point := line2.Points[0]
			if Distance(point, lastPoint) < joinFraction*(l.Bounds[2]+l.Bounds[3])/2 {
				l2.Lines[len(l2.Lines)-1].Points = append(l2.Lines[len(l2.Lines)-1].Points, line2.Points...)
			} else {
				l2.Lines = append(l2.Lines, line2)
			}
		} else {
			l2.Lines = append(l2.Lines, line2)
		}
	}
	log.Debugf("consolidated %d lines to %d lines", len(l.Lines), len(l2.Lines))
	return
}

func (l *Line) Add(x, y float64) {
	l.Points = append(l.Points, NewPoint(x, y))
}

func (l *Line) Length() int {
	return len(l.Points)
}

func NewPoint(x, y float64) Point {
	return Point{x, y}
}

func NewLine() Line {
	return Line{[]Point{}}
}
