package systems

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"fmt"
	"image/color"
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
	CellMap [][]int
}

// This type represents an cell (entity), this cell
// has the BasicEntity (standard unique id for entities)
// and two components, render and space(transform)
type MapCell struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	CellComponent
}

const (
	NORMAL = iota
	GOAL
	NEST
	OBSTACLE
	ANT
	PHEROMONE
)

// This type is a component to be added to cells
// represents some info
type CellComponent struct {
	ants      int
	pheromone float32
	cellType  int
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

	for _, cell := range mcs.CellMap {
		fmt.Println(cell)
	}

	// Calculate cell size
	size := engo.WindowWidth() / float32(len(mcs.CellMap))
	size *= 0.9

	for x := range mcs.CellMap {
		for y := range mcs.CellMap[x] {
			//calculate x and y depending on number of cells
			newCell(ecsWorld, engo.Point{(engo.WindowWidth() / (2 * float32(len(mcs.CellMap)))) * float32((x*2 + 1)), (engo.WindowHeight() / (2 * float32(len(mcs.CellMap[x])))) * float32((y*2 + 1))}, size, mcs.CellMap[x][y], x, y)
		}
	}
}

func newCell(world *ecs.World, position engo.Point, size float32, cellType int, x int, y int) {
	// Entities must be initiated
	cell := MapCell{BasicEntity: ecs.NewBasic()}

	// Setting up space component
	cell.SpaceComponent = common.SpaceComponent{
		Position: position,
		Width:    size,
		Height:   size,
	}

	// Setting up Render Component
	cell.RenderComponent = common.RenderComponent{
		Drawable: common.Rectangle{},
		Scale:    engo.Point{0.98, 0.98},
		Color:    color.White,
	}

	if cellType != 0 {
		fmt.Println("Found the goal:", cellType)
	}

	cell.CellComponent = CellComponent{
		ants:      0,
		pheromone: 0,
		cellType:  cellType,
	}

	// Adding the entity to the RenderSystem
	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&cell.BasicEntity, &cell.RenderComponent, &cell.SpaceComponent)
		case *PainterSystem:
			sys.Add(&cell.BasicEntity, &cell.RenderComponent, &cell.CellComponent, &cell.SpaceComponent)

		case *AntMoverSystem:
			sys.Add(&cell.BasicEntity, &cell.RenderComponent, &cell.CellComponent, &cell.SpaceComponent, x, y)
		}
	}

}
