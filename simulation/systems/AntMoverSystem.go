package systems

import (
	"engo.io/ecs"
	"engo.io/engo/common"
	"fmt"
	"image/color"
)

type AntMoverSystem struct {
	cells []cellEntity
}

type cellEntity struct {
	ecs.BasicEntity
	*common.RenderComponent
	common.SpaceComponent
	CellComponent
}

// Called when entity is removed from the world so
// that it can also be removed from this system
func (*PainterSystem) Remove(ecs.BasicEntity) {
}

func (ps *PainterSystem) Update(dt float32) {
}

// Called when system is initialized
func (ps *PainterSystem) New(world *ecs.World) {

}

// Called to add entities to the system
func (ps *PainterSystem) Add(basic *ecs.BasicEntity, render *common.RenderComponent, cell *CellComponent, space *common.SpaceComponent) {
	fmt.Println("Added cell:", cell.cellType)
	ps.cells = append(ps.cells, cellEntity{*basic, render, *space, *cell})
}
