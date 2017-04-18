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
	ExportMap() ([][]int, map[int][2]int)
	ExportRoutes() [][][]int
	ExportPheromones() [][][]float64
}

//AntColonyOptimization struct for the algorithm
type AntColonyOptimization struct {
	world           world.World
	size            int
	antArray        []*ants.Ant
	iterations      int
	routeMatrix     [][][]int
	pheromoneMatrix [][][]float64
}

//Init initializes the ACO
func (aco *AntColonyOptimization) Init(kants, size, iterations int) {

	//WORLD
	var impl = new(world.WorldImpl)
	var w world.World = impl

	//ANT
	var antArray = make([]*ants.Ant, kants)
	aco.size = size
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

	//route matrix init
	aco.routeMatrix = make([][][]int, kants)
	for i := range aco.routeMatrix {
		aco.routeMatrix[i] = make([][]int, iterations)
		for j := range aco.routeMatrix[i] {
			aco.routeMatrix[i][j] = make([]int, 0)
		}
	}
	//phe matrix init
	aco.pheromoneMatrix = make([][][]float64, iterations)
	for i := range aco.pheromoneMatrix {
		aco.pheromoneMatrix[i] = make([][]float64, size)
		for j := range aco.pheromoneMatrix[i] {
			aco.pheromoneMatrix[i][j] = make([]float64, size)
		}
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
	for iter := 0; iter < aco.iterations; iter++ {

		for _, ant := range aco.antArray {
			ant.Reset()
			//fmt.Println(ant)
			go ant.Run(cha)

		}
		//wait kAnts times receive channel
		for k := 0; k < len(aco.antArray); k++ {
			route := <-cha
			//var intRoute []int
			for i, val := range route {
				//intRoute = append(intRoute, int(val))
				aco.routeMatrix[k][iter] = append(aco.routeMatrix[k][iter], int(val))
				if i+1 < len(route) {
					aco.world.PutPheromone(route[i], route[i+1])
				}
			}

		}
		//pheromone snapshot
		aco.world.UpdatePheromones()
		snapshot := aco.world.ExportPheromones()
		for i := 0; i < aco.size; i++ {
			for j := 0; j < aco.size; j++ {
				aco.pheromoneMatrix[iter][i][j] = snapshot[i][j]
			}
		}
	}

}

//ExportMap exports the map used by the algorithm to use in the simulation
func (aco AntColonyOptimization) ExportMap() ([][]int, map[int][2]int) {
	return aco.world.ExportMap()
}

//ExportRoutes exports a matrix of routes for the simulator to play
func (aco AntColonyOptimization) ExportRoutes() [][][]int {
	//order of routes is cagative for simulation
	return aco.routeMatrix
}

//ExportPheromones dfa
func (aco AntColonyOptimization) ExportPheromones() [][][]float64 {
	return aco.pheromoneMatrix
}
