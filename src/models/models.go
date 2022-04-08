package models

type Lines struct {
	Lines  []Line
	Height float64
	Width  float64
}

type Line struct {
	Points []Point
}

type Point struct {
	X, Y float64
}

func (ls *Lines) Add(l Line) {
	ls.Lines = append(ls.Lines, l)
	// TODO: recompute the height and width
}

func (l Lines) BoundingBox(x, y, w, h float64) (l2 Lines) {
	l2 = l.Normalize()

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
	minPoint := Point{0, 0}
	for _, line := range l2.Lines {
		for _, p := range line.Points {
			if p.X < minPoint.X {
				minPoint.X = p.X
			}
			if p.Y < minPoint.Y {
				minPoint.Y = p.Y
			}
		}
	}
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
