package grad

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestGA(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	points := []Point{
		{2, 4},
		Point{1, 1},
		Point{1.1, 5},
		Point{3, 1},
		Point{4, 2},
		Point{5, 5},
		Point{7, 1},
		Point{3, 1.5},
	}
	path := FindBestPath(points)
	fmt.Println(path)
}
