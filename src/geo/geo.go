package geo

import (
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"os"
	"runtime"
	"sort"

	"github.com/fogleman/gg"
	log "github.com/schollz/logger"
	"github.com/schollz/progressbar/v3"
)

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

func (l Lines) Copy() (l2 Lines) {
	l2 = Lines{}
	for _, line := range l.Lines {
		line2 := Line{}
		for _, p := range line.Points {
			line2.Points = append(line2.Points, Point{p.X, p.Y})
		}
		l2.Lines = append(l2.Lines, line2)
	}
	return

}

// Simplify points by some percentage using the
// https://en.wikipedia.org/wiki/Visvalingam%E2%80%93Whyatt_algorithm
func (l Lines) Simplify(percent float64) (l2 Lines) {

	pointsAvailable := 0.0
	for _, line := range l.Lines {
		if len(line.Points) > 2 {
			pointsAvailable += float64(len(line.Points) - 2)
		}
	}

	log.Tracef("found %2.0f available points to simplify", pointsAvailable)

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

	l2 = Lines{}
	for linei, line := range l.Lines {
		line2 := Line{}
		for pointi, p := range line.Points {
			if _, ok := removal[fmt.Sprintf("%d-%d", linei, pointi)]; !ok {
				line2.Points = append(line2.Points, Point{p.X, p.Y})
			}
		}
		l2.Lines = append(l2.Lines, line2)
	}

	return
}

func (l *Lines) Add(line Line) {
	l.Lines = append(l.Lines, line)
}

func (l Lines) Draw(fname string) (err error) {
	im := l.DrawStep(-1)
	f, err := os.Create(fname)
	if err != nil {
		return
	}
	defer f.Close()

	err = png.Encode(f, im)
	if err != nil {
		return
	}
	return
}

func (l Lines) DrawStep(step int, im0 ...image.Image) (im image.Image) {
	var dc *gg.Context
	if len(im0) > 0 {
		dc = gg.NewContextForImage(im0[0])
	} else {
		dc = gg.NewContext(int(l.Bounds[0]+l.Bounds[2]), int(l.Bounds[1]+l.Bounds[3]))
		dc.DrawRectangle(0, 0, l.Bounds[0]+l.Bounds[2], l.Bounds[1]+l.Bounds[3])
		dc.SetRGB(1, 1, 1)
		dc.Fill()
	}
	i := 0
	for _, line := range l.Lines {
		for j, point := range line.Points {
			if j == 0 {
				continue
			}
			dc.DrawLine(line.Points[j-1].X, line.Points[j-1].Y, point.X, point.Y)
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
	dc.SetLineWidth(4)
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
		outGif.Delay = append(outGif.Delay, 0)
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
