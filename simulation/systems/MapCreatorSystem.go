package systems

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"fmt"
	"github.com/pedroallenrevez/goAnt/world"
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

type MapCreatorSystem struct {
	world   *ecs.World
	cellMap [][]int
}

// This type represents an cell (entity), this cell
// has the BasicEntity (standard unique id for entities)
// and two components, render and space(transform)
type MapCell struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

// Called when entity is removed from the world so
// that it can also be removed from this system
func (*MapCreatorSystem) Remove(ecs.BasicEntity) {}

func (mcs *MapCreatorSystem) Update(dt float32) {

}

// Called when system is initialized
func (mcs *MapCreatorSystem) New(ecsWorld *ecs.World) {
	fmt.Println("MapCreatorSystem INITIALIZED")

	// Save world reference to use it on update to add
	// Entity to world
	mcs.world = ecsWorld

	// Import map from AntColony
	antMap := world.WorldImpl{}
	antMap.Init(5, func(world.Point, world.Point) float64 { return 1 })
	fmt.Println("World map INITIALIZED")
	mcs.cellMap = antMap.ExportMap()

	for _, cell := range mcs.cellMap {
		fmt.Println(cell)
	}

	// Calculate cell size
	size := engo.WindowWidth() / float32(len(mcs.cellMap))
	size *= 0.9

	for x := range mcs.cellMap {
		for y := range mcs.cellMap[x] {
			//calculate x and y depending on number of cells
			newCell(ecsWorld, engo.Point{(engo.WindowWidth() / (2 * float32(len(mcs.cellMap)))) * float32((x*2 + 1)), (engo.WindowHeight() / (2 * float32(len(mcs.cellMap[x])))) * float32((y*2 + 1))}, size)
		}
	}
}

func newCell(world *ecs.World, position engo.Point, size float32) {
	// Entities must be initiated
	cell := MapCell{BasicEntity: ecs.NewBasic()}

	// Setting up space component
	cell.SpaceComponent = common.SpaceComponent{
		Position: position,
		Width:    size,
		Height:   size,
	}

	// Setting up Render Component

	texture, err := common.LoadedSprite("textures/Ant.png")

	if err != nil {
		fmt.Println("Unable to load texture: " + err.Error())
	}

	cell.RenderComponent = common.RenderComponent{
		Drawable: texture,
		Scale:    engo.Point{0.5, 0.5},
	}

	// Adding the entity to the RenderSystem
	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&cell.BasicEntity, &cell.RenderComponent, &cell.SpaceComponent)
		}
	}

}
