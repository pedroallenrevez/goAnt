package main

import (
	"fmt"
	"world"
)

func main() {
	world := WorldImpl{}
	world.init()
	//define distance calculation function (maybe include default in world?)
	//init k ants(threads)
	//start ant threads
	//join threads
	world.updatePheromones()
	//updateWorld with ant routes ( shortest dist, pheromone) <-- Pheromone updating is not necessary, what do you mean with updateWorld with shortest dist?
	//^ if this refers to keeping distance tables to allow faster routing, it's not implemented
	//ant reset
}
