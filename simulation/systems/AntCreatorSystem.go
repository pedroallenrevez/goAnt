package systems

import (
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

type AntCreatorSystem struct{}

// Called when entity is removed from the world so
// that it can also be removed from this system
func (*AntCreatorSystem) Remove(ecs.BasicEntity) {}

func (*AntCreatorSystem) Update(dt float32) {}

func (*AntCreatorSystem) New(world *ecs.World) {
	fmt.Println("AntCreatorSystem INITIALIZED")
}
