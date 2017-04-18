package main

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"fmt"
	"github.com/pedroallenrevez/goAnt/aco"
	"github.com/pedroallenrevez/goAnt/simulation/systems"
	"image/color"
	"os"
	"strconv"
)

// BASICS:
// Each scene contains 1 world, multiple systems
// and a multitude of entities
type myScene struct{}

// This type represents an ant (entity), this ant
// has the BasicEntity (standard unique id for entities)
// and two components, render and space(transform)
type Ant struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

// Type uniquely identifies the game type
func (*myScene) Type() string {
	return "myGame"
}

// Preload is called before loading any assets
// allows registering/queueing them
func (*myScene) Preload() {
	// Assumes to be inside assets folder
	engo.Files.Load("textures/Ant.png")
}

// Setup is called before main loop starts.
// This is where you add entitites and systems
// to the scene
func (sc *myScene) Setup(world *ecs.World) {

	// Input needs to be registered
	// engo.Input.RegisterButton("AddAnt", engo.F1)

	common.SetBackground(color.Black)

	//seconds each round lasts in AntMover
	seconds, _ := strconv.ParseFloat(os.Args[4], 32)

	//Initialize ACO script
	acoInstance := aco.AntColonyOptimization{}
	antNum, _ := strconv.Atoi(os.Args[1])
	mapSize, _ := strconv.Atoi(os.Args[2])
	iterNum, _ := strconv.Atoi(os.Args[3])

	fmt.Println("STARTING WITH ANTS:", antNum, " MAP:", mapSize, " ITER:", iterNum)

	acoInstance.Init(antNum, mapSize, iterNum)
	//cellMap, mapToMap := acoInstance.ExportMap()
	cellMap, mapToMap := acoInstance.ExportMap()

	acoInstance.Run()
	antRoutes := acoInstance.ExportRoutes()

	for _, ant := range antRoutes {
		for iteration := range ant {

			ant[iteration] = append([]int{0}, ant[iteration]...)
			ant[iteration] = append(ant[iteration], reverseArray(ant[iteration])[1:]...)
		}
	}
	// Systems need to be added to the world
	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&common.MouseSystem{})

	// Initialize custom systems last to make sure their
	// depencies are already initialized
	// world.AddSystem(&systems.AntCreatorSystem{})
	world.AddSystem(&systems.PainterSystem{})
	world.AddSystem(&systems.AntMoverSystem{Routes: antRoutes, IDtoMap: mapToMap, SecondsPerRound: float32(seconds)})
	world.AddSystem(&systems.MapCreatorSystem{CellMap: cellMap})

}

func main() {
	opts := engo.RunOptions{
		Title:  "Hello World",
		Width:  800,
		Height: 800,
	}

	engo.Run(opts, &myScene{})

}

func reverseArray(route []int) []int {
	result := make([]int, len(route))
	for i, j := len(route)-1, 0; i >= 0; i, j = i-1, j+1 {
		result[j] = route[i]
	}
	fmt.Println("Reversed:", route, " into:", result)
	return result
}
