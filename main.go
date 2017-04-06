package main

import (
	"fmt"
	"github.com/pedroallenrevez/goAnt/ants"
	"github.com/pedroallenrevez/goAnt/world"
)

func main() {
	//add argv support
	//numbers of ants
	kants := 50
	ants := make([]Ant, kants)
	//init world
	/*
	   var a IFace
	   a = Create()
	   a.SetSomeField("World")
	   fmt.Println(a.GetSomeField())
	*/
	var world World
	worldImpl := WorldImpl{}
	world.init()
	world = worldImpl
	fmt.Println("Begin ant creation")
	for i = 0; i < kants; i++ {
		// new -> allocates a zero value of type T and returns pointer
		newAnt := new(Ant)
		newAnt.init(world)
		append(ants, newAnt)
	}

	euclidean := func(start, end Point) float64 {
		return math.Pow(start.x-end.x, 2) + math.Pow(start.y-end.y, 2)
	}
	//define distance calculation function (maybe include default in world?)
	//init k ants(threads)
	//start ant threads
	//join threads
	//world.updatePheromones()
	//updateWorld with ant routes ( shortest dist, pheromone) <-- Pheromone updating is not necessary, what do you mean with updateWorld with shortest dist?
	//^ if this refers to keeping distance tables to allow faster routing, it's not implemented
	//ant reset
}
