package grad

import (
	"math"
	"time"

	"github.com/jmcvetta/randutil"
	log "github.com/schollz/logger"
)

type Point struct {
	X, Y float64
}

func FindBestPath(points []Point) (bestPath []int) {
	bestPathLength := 100000000.0
	currentTime := time.Now()
	for tries := 0; tries < 10000000; tries++ {
		if time.Since(currentTime) > time.Duration(1*time.Second) {
			break
		}
		for i := 0; i < len(points); i++ {
			path := FindPath(points, i)
			pathLength := PathLength(path, points)
			if pathLength < bestPathLength {
				log.Debugf("best path: %+v: %2.3f", path, pathLength)
				bestPathLength = pathLength
				bestPath = path
			}
		}
	}
	return
}

func PathLength(path []int, points []Point) (pathLength float64) {
	for i := 2; i < len(points); i += 2 {
		otherPoint := points[path[i-1]]
		point := points[path[i]]
		pathLength += math.Sqrt(math.Pow(point.X-otherPoint.X, 2) + math.Pow(point.Y-otherPoint.Y, 2))
	}
	return
}

func FindPath(points []Point, start int) (path []int) {
	path = []int{start}
	for {
		if path[len(path)-1]%2 == 0 {
			path = append(path, path[len(path)-1]+1)
		} else {
			path = append(path, path[len(path)-1]-1)
		}

		if len(path) == len(points) {
			break
		}
		finished := make(map[int]struct{})
		for _, v := range path {
			finished[v] = struct{}{}
		}

		// find next closest point
		path = append(path, ClosestPoint(points[path[len(path)-1]], points, finished))
	}
	return
}

type Dist struct {
	Val float64
	Ind int
}

func ClosestPoint(point Point, points []Point, finished map[int]struct{}) (besti int) {
	// find closest point based on weighted randomness
	choices := []randutil.Choice{}
	for i, otherPoint := range points {
		if _, ok := finished[i]; ok {
			continue
		}
		weight := 10000000
		dist := math.Sqrt(math.Pow(point.X-otherPoint.X, 2) + math.Pow(point.Y-otherPoint.Y, 2))
		if dist > 0 {
			weight = int(float64(weight) / dist)
		}
		choices = append(choices, randutil.Choice{
			Weight: weight,
			Item:   i,
		})
	}
	choice, _ := randutil.WeightedChoice(choices)
	besti = choice.Item.(int)
	return
}
