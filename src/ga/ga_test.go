package ga

import (
	"math/rand"
	"testing"
	"time"
)

func TestGA(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	cities := []City{
		City{2, 4},
		City{1, 1},
		City{1.1, 5},
		City{3, 1},
		City{4, 2},
		City{5, 5},
		City{7, 1},
		City{3, 1.5},
	}
	trips := []Trip{
		Trip{[]int{0}},
		Trip{[]int{1, 2}},
		Trip{[]int{3, 4}},
		Trip{[]int{5}},
		Trip{[]int{6}},
		Trip{[]int{7}},
	}
	w := NewWorld(cities, trips)
	// j := w.RandomJourney()
	// j.Plot("1.png")
	// j.Mutate(1)
	// j.Plot("2.png")

	j := w.FindBest()
	j.Plot("1.png")

}
