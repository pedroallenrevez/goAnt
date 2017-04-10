package main

import (
	//"fmt"
	//"github.com/pedroallenrevez/goAnt/ants"
	//"github.com/pedroallenrevez/goAnt/world"
	//"math"
	"github.com/pedroallenrevez/goAnt/aco"
)

/*
//WORLD
var impl = new(world.WorldImpl)
var w world.World = impl

//ANT
var kants = 1
var iterations = 10
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
	w.Init(20, euclidean)
	//fmt.Println(w)

	//init k ants(threads)
	for i := 0; i < kants; i++ {
		newAnt := ants.NewAnt(w)
		ant = newAnt
		antArray[i] = newAnt
	}

	//start ant threads
	cha := make(chan []world.NodeID)
	//reset and run ants
	fmt.Println("running ants")
	for iter := 0; iter <= iterations; iter++ {

		for _, ant := range antArray {
			ant.Reset()
			//fmt.Println(ant)
			go ant.Run(cha)

		}
		//wait kAnts times receive channel
		for i := 0; i < len(antArray); i++ {
			route := <-cha

			for i := range route {
				if i+1 < len(route) {
					w.PutPheromone(route[i], route[i+1])
				}
			}

		}
		w.UpdatePheromones()
	}

}
*/

func main() {
	var impl = new(aco.AntColonyOptimization)
	var aco aco.ACOInterface = impl
	//kants, size, iter
	aco.Init(5, 20, 10)
	aco.Run()
}
