package main

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
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
func (*myScene) Setup(world *ecs.World) {

	// Input needs to be registered
	engo.Input.RegisterButton("AddAnt", engo.F1)

	common.SetBackground(color.White)

	// Systems need to be added to the world
	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&common.MouseSystem{})

	// Initialize custom systems last to make sure their
	// depencies are already initialized
	world.AddSystem(&systems.AntCreatorSystem{})

	/*
		// Entities must be initiated
		ant := Ant{BasicEntity: ecs.NewBasic()}

		// Setting up space component
		ant.SpaceComponent = common.SpaceComponent{
			Position: engo.Point{10, 10},
			Width:    303,
			Height:   641,
		}

		// Setting up Render Component

		texture, err := common.LoadedSprite("textures/Ant.png")

		if err != nil {
			fmt.Println("Unable to load texture: " + err.Error())
		}

		ant.RenderComponent = common.RenderComponent{
			Drawable: texture,
			Scale:    engo.Point{1, 1},
		}

		// Adding the entity to the RenderSystem
		for _, system := range world.Systems() {
			switch sys := system.(type) {
			case *common.RenderSystem:
				sys.Add(&ant.BasicEntity, &ant.RenderComponent, &ant.SpaceComponent)
			}
		}
	*/

}

func main() {
	opts := engo.RunOptions{
		Title:  "Hello World",
		Width:  400,
		Height: 400,
	}

	engo.Run(opts, &myScene{})

}
