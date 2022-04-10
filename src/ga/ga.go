package ga

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"

	log "github.com/schollz/logger"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func Distance(a City, b City) float64 {
	return math.Sqrt(math.Pow(a.X-b.X, 2) + math.Pow(a.Y-b.Y, 2))
}

// City is a location in the world
type City struct {
	X, Y float64
}

// World is the list of all the cities
// and the possible trips
type World struct {
	Cities []City
	Trips  []Trip
}

// Trip can be one city or two cities in a specific order
type Trip struct {
	Dests []int
}

// Journey is the total journey
type Journey struct {
	Cities []City
	Trips  []Trip
	Length float64
}

func NewWorld(cities []City, trips []Trip) (w World) {
	w = World{Cities: cities, Trips: trips}
	return
}

func (w World) RandomJourney() (j Journey) {
	trips := []Trip{}
	for _, t := range w.Trips {
		trips = append(trips, t)
	}

	for i := range trips {
		j := rand.Intn(i + 1)
		trips[i], trips[j] = trips[j], trips[i]
	}

	j = Journey{Cities: w.Cities, Trips: trips}
	j.Eval()
	return
}

func (j *Journey) Eval() {
	pts := make(plotter.XYs, len(j.Cities))
	i := 0
	for _, trip := range j.Trips {
		for _, dest := range trip.Dests {
			city := j.Cities[dest]
			pts[i].X = city.X
			pts[i].Y = city.Y
			i++
		}
	}

	dist := 0.0
	for i, pt := range pts {
		if i == 0 {
			continue
		}
		dist += math.Sqrt(math.Pow(pt.X-pts[i-1].X, 2) + math.Pow(pt.Y-pts[i-1].Y, 2))
	}
	j.Length = dist
}

func (j Journey) Plot(fname string) (err error) {
	rand.Seed(int64(0))

	p := plot.New()

	p.Title.Text = fmt.Sprintf("%2.2f", j.Length)
	p.X.Label.Text = ""
	p.Y.Label.Text = ""
	p.HideAxes()

	pts := make(plotter.XYs, len(j.Cities))
	i := 0
	for _, trip := range j.Trips {
		for _, dest := range trip.Dests {
			city := j.Cities[dest]
			pts[i].X = city.X
			pts[i].Y = city.Y
			i++
		}
	}
	err = plotutil.AddLinePoints(p,
		"", pts,
	)
	if err != nil {
		log.Error(err)
		return
	}

	// Save the plot to a PNG file.
	if err = p.Save(3*vg.Inch, 3*vg.Inch, fname); err != nil {
		log.Error(err)
	}
	return
}

func (jo *Journey) Mutate(n int) {
	rand.Seed(time.Now().UnixNano())
	for times := 0; times < n; times++ {
		i := rand.Intn(len(jo.Trips))
		j := rand.Intn(len(jo.Trips))
		jo.Trips[i], jo.Trips[j] = jo.Trips[j], jo.Trips[i]
	}
	jo.Eval()
}

// func (j0 Journey) Crossover(k Journey) Journey {
// 	jnew := Journey{
// 		Cities: j0.Cities,
// 		Trips:  j0.Trips,
// 	}
// 	tripNext := make[map[string]string]
// 	for i := range jnew.Trips {
// 		j := rand.Intn(i + 1)
// 		jnew.Trips[i], jnew.Trips[j] = jnew.Trips[j], jnew.Trips[i]
// 	}
// 	// remove duplicates
// 	trips := make([]Trip, len(j0.Trips))
// 	tripsExist := make(map[string]struct{})
// 	i := 0
// 	for _, trip := range jnew.Trips {
// 		tripName := fmt.Sprintf("%+v", trip)
// 		if _, ok := tripsExist[tripName]; ok {
// 			continue
// 		}
// 		tripsExist[tripName] = struct{}{}
// 		trips[i] = trip
// 		i++
// 	}
// 	jnew.Trips = trips
// 	jnew.Eval()
// 	return jnew
// }

type GA struct {
	MutationRate float64
	Population   int
	Journies     []Journey
}

func NewGA(w World, population int, mutationRate float64) GA {
	ga := GA{
		MutationRate: mutationRate,
		Population:   population,
		Journies:     make([]Journey, population),
	}
	for i := 0; i < ga.Population; i++ {
		ga.Journies[i] = w.RandomJourney()
	}
	return ga
}

func (ga *GA) Iterate() {
	sort.Slice(ga.Journies, func(i, j int) bool {
		return ga.Journies[i].Length < ga.Journies[j].Length
	})
	journies := make([]Journey, ga.Population)
	for i := 0; i < ga.Population; i++ {
		trips := make([]Trip, len(ga.Journies[0].Trips))
		for j, trip := range ga.Journies[rand.Intn(10)].Trips {
			trips[j] = Trip{trip.Dests}
		}
		journies[i] = Journey{
			Trips:  trips,
			Cities: ga.Journies[0].Cities,
		}
		//journies[i] = journies[i].Crossover(ga.Journies[rand.Intn(4)])
		if i > 1 && rand.Float32() < float32(ga.MutationRate) {
			journies[i].Mutate(rand.Intn(10) + 1)
		}
		journies[i].Eval()
	}
	for i := 0; i < ga.Population; i++ {
		ga.Journies[i] = journies[i]
	}
}

func (w World) FindBest() Journey {
	ga := NewGA(w, 200, 0.9)
	last := 100000.0
	j := 0
	for i := 0; i < 500; i++ {
		ga.Iterate()
		if ga.Journies[0].Length < last {
			// ga.Journies[0].Plot(fmt.Sprintf("%06d.png", j))
			j++
		}
		last = ga.Journies[0].Length
	}
	_ = j
	return ga.Journies[0]
}
