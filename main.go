package main

import (
	"fmt"
	"github.com/pedroallenrevez/goAnt/ants"
	"github.com/pedroallenrevez/goAnt/world"
	"math"
)

//WORLD
var impl = new(world.WorldImpl)
var w world.World = impl

//ANT
var kants = 50
var antArray = make([]*ants.Ant, kants)
var ant ants.AntInterface

//EXPORT MATRIX

func main() {
	//distance callback
	euclidean := func(start, end world.Point) float64 {
		distanceX := math.Abs(float64(start.X - end.X))
		distanceY := math.Abs(float64(start.Y - end.Y))
		return math.Pow(distanceX, 2) + math.Pow(distanceY, 2)
	}
	w.Init(100, euclidean)

	//init k ants(threads)
	for i := 0; i < kants; i++ {
		newAnt := ants.NewAnt(w)
		ant = newAnt
		antArray = append(antArray, newAnt)
	}

	//start ant threads
	cha := make(chan []int)
	//reset and run ants
	fmt.Println("running ants")
	for _, ant := range antArray {
		//ant.Init(w)
		ant.Run(cha)

	}
	for i := 0; i < len(antArray); i++ {
		//wait kAnts times receive channel
		route := <-cha
		//recive from channel and add to map
		//update pheromonemap with putpheromone
		for i := range route {
			if i+1 <= len(route) {
				w.PutPheromone(world.NodeID(route[i]), world.NodeID(route[i+1]))
			}

		}

	}
	//updateWorld with ant routes ( shortest dist, pheromone) <-- Pheromone updating is not necessary, what do you mean with updateWorld with shortest dist?
	//w.UpdatePheromones()

}
