package main

import (
	"fmt"
	"github.com/pedroallenrevez/goAnt/aco"
)

func main() {
	var impl = new(aco.AntColonyOptimization)
	var aco aco.ACOInterface = impl
	//kants, size, iter
	aco.Init(1, 5, 2)
	aco.Run()
	mapa, table := aco.ExportMap()

	fmt.Println(mapa, table)
	routes := aco.ExportRoutes()
	fmt.Println(routes)
}
