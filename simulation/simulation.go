package main

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"github.com/pedroallenrevez/goAnt/aco"
	"github.com/pedroallenrevez/goAnt/simulation/systems"
	"image/color"
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

	//Initialize ACO script
	acoInstance := aco.AntColonyOptimization{}
	acoInstance.Init(2, 10, 2)
	//cellMap, mapToMap := acoInstance.ExportMap()
	cellMap, _ := acoInstance.ExportMap()

	// Systems need to be added to the world
	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&common.MouseSystem{})

	// Initialize custom systems last to make sure their
	// depencies are already initialized
	// world.AddSystem(&systems.AntCreatorSystem{})
	world.AddSystem(&systems.PainterSystem{})
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
