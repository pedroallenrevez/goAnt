package systems

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"fmt"
)

/*
// System is an interface that implements an ECS-System
// It should iterate over its entitites on `Update` and
// do what it is suitable on the current implementation
type System interface {
	// Run every frame, dt = frame delta time
	Update(dt float32)

	// Initialisation of the System
	New(*ecs.World)

	// Removes entity from the system
	Remove(ecs.BasicEntity)
}
*/

type MouseTracker struct {
	ecs.BasicEntity
	common.MouseComponent
}

type AntCreatorSystem struct {
	world        *ecs.World
	mouseTracker MouseTracker
}

// This type represents an ant (entity), this ant
// has the BasicEntity (standard unique id for entities)
// and two components, render and space(transform)
type Ant struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

// Called when entity is removed from the world so
// that it can also be removed from this system
func (*AntCreatorSystem) Remove(ecs.BasicEntity) {}

func (acs *AntCreatorSystem) Update(dt float32) {
	if engo.Input.Button("AddAnt").JustPressed() {

		// Entities must be initiated
		ant := Ant{BasicEntity: ecs.NewBasic()}

		// Setting up space component
		ant.SpaceComponent = common.SpaceComponent{
			Position: engo.Point{acs.mouseTracker.MouseX, acs.mouseTracker.MouseY},
			Width:    30,
			Height:   30,
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
		for _, system := range acs.world.Systems() {
			switch sys := system.(type) {
			case *common.RenderSystem:
				sys.Add(&ant.BasicEntity, &ant.RenderComponent, &ant.SpaceComponent)
			}
		}
	}
}

// Called when system is initialized
func (acs *AntCreatorSystem) New(world *ecs.World) {
	fmt.Println("AntCreatorSystem INITIALIZED")

	// Save world reference to use it on update to add
	// Entity to world
	acs.world = world

	acs.mouseTracker.BasicEntity = ecs.NewBasic()

	// Track: true makes it so that you can always know the mouse Position
	// not only when it is interacting with the entity
	acs.mouseTracker.MouseComponent = common.MouseComponent{Track: true}

	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.MouseSystem:
			sys.Add(&acs.mouseTracker.BasicEntity, &acs.mouseTracker.MouseComponent, nil, nil)
		}
	}

}
