package aco

import (
	"fmt"
	"github.com/pedroallenrevez/goAnt/ants"
	"github.com/pedroallenrevez/goAnt/world"
	"math"
)

//ACOInterface defines the interface to use the ACO algorithm
type ACOInterface interface {
	//kants, size, iterations
	Init(kants, size, iterations int)
	Run()
	ExportMap()
	ExportRoutes()
}

//AntColonyOptimization struct for the algorithm
type AntColonyOptimization struct {
	world      world.World
	antArray   []*ants.Ant
	iterations int
}

//Init initializes the ACO
func (aco *AntColonyOptimization) Init(kants, size, iterations int) {

	//WORLD
	var impl = new(world.WorldImpl)
	var w world.World = impl

	//ANT
	var antArray = make([]*ants.Ant, kants)

	//distance callback
	euclidean := func(start, end world.Point) float64 {
		distanceX := math.Abs(float64(start.X - end.X))
		distanceY := math.Abs(float64(start.Y - end.Y))
		return math.Pow(distanceX, 2) + math.Pow(distanceY, 2)
	}
	w.Init(size, euclidean)
	//fmt.Println(w)

	//init k ants(threads)
	for i := 0; i < kants; i++ {
		newAnt := ants.NewAnt(w)
		antArray[i] = newAnt
	}
	aco.antArray = antArray
	aco.world = w
	aco.iterations = iterations
}

//Run runs the algorithm with specified parameters
func (aco *AntColonyOptimization) Run() {
	cha := make(chan []world.NodeID)
	//reset and run ants
	fmt.Println("running ants")
	for iter := 0; iter <= aco.iterations; iter++ {

		for _, ant := range aco.antArray {
			ant.Reset()
			//fmt.Println(ant)
			go ant.Run(cha)

		}
		//wait kAnts times receive channel
		for k := 0; k < len(aco.antArray); k++ {
			route := <-cha
			for i := range route {
				if i+1 < len(route) {
					aco.world.PutPheromone(route[i], route[i+1])
				}
			}

		}
		aco.world.UpdatePheromones()
	}
}

//ExportMap exports the map used by the algorithm to use in the simulation
func (aco AntColonyOptimization) ExportMap() {}

//ExportRoutes exports a matrix of routes for the simulator to play
func (aco AntColonyOptimization) ExportRoutes() {}
