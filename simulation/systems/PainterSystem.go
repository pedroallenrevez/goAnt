package systems

import (
	"engo.io/ecs"
	"engo.io/engo/common"
	"image/color"
)

type PainterSystem struct {
	Cells []cellEntity
}

type cellEntity struct {
	ecs.BasicEntity
	*common.RenderComponent
	*common.SpaceComponent
	*CellComponent
}

// Called when entity is removed from the world so
// that it can also be removed from this system
func (*PainterSystem) Remove(ecs.BasicEntity) {}

func (ps *PainterSystem) Update(dt float32) {

	var cellColor color.Color

	//TODO, why color all cells? keep status and only color the ones who change

	for _, cell := range ps.Cells {
		if cell.ants > 0 {
			cellColor = color.RGBA{0, 255, 0, 255}
		} else if cell.cellType == NORMAL {
			cellColor = color.RGBA{255, 255, 255, 255}
		} else if cell.cellType == GOAL {
			cellColor = color.RGBA{255, 0, 0, 255}
		} else if cell.cellType == NEST {
			cellColor = color.RGBA{0, 255, 0, 255}
		} else if cell.cellType == OBSTACLE {
			cellColor = color.RGBA{0, 0, 255, 255}
		}

		cell.RenderComponent.Color = cellColor
	}
}

// Called when system is initialized
func (ps *PainterSystem) New(world *ecs.World) {

}

// Called to add entities to the system
func (ps *PainterSystem) Add(basic *ecs.BasicEntity, render *common.RenderComponent, cell *CellComponent, space *common.SpaceComponent) {
	ps.Cells = append(ps.Cells, cellEntity{*basic, render, space, cell})
}
