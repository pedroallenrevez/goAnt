package systems

import (
	"engo.io/ecs"
	"engo.io/engo/common"
	"fmt"
)

type AntMoverSystem struct {
	cells           map[coordinate]cellEntity
	Routes          [][][]int
	IDtoMap         map[int][2]int
	SecondsPerRound float32
}

type coordinate struct {
	x int
	y int
}

var roundTimer float32
var roundCounter int
var iterationCounter int
var over bool
var pheroActivate []bool

// Called when entity is removed from the world so
// that it can also be removed from this system
func (*AntMoverSystem) Remove(ecs.BasicEntity) {
}

func (ams *AntMoverSystem) Update(dt float32) {
	updateRoundTimer(dt)
	if !over {
		if roundTimer >= ams.SecondsPerRound {
			roundTimer = 0
			updates := 0
			fmt.Println("Iteration:", iterationCounter, " Round:", roundCounter)

			for ant := range ams.Routes {
				if roundCounter < len(ams.Routes[ant][iterationCounter]) {
					coords := ams.IDtoMap[ams.Routes[ant][iterationCounter][roundCounter]]
					ams.cells[coordinate{x: coords[0], y: coords[1]}].ants += 1
					if ams.cells[coordinate{coords[0], coords[1]}].cellType == 1 {
						fmt.Println("ANT:", ant, " found GOAL! at:", coords[0], ",", coords[1])
						pheroActivate[ant] = true
					}
					//Delete ants from cells they left
					if roundCounter > 0 {
						oldCoords := ams.IDtoMap[ams.Routes[ant][iterationCounter][roundCounter-1]]
						ams.cells[coordinate{x: oldCoords[0], y: oldCoords[1]}].ants -= 1

						if pheroActivate[ant] == true {
							ams.cells[coordinate{x: oldCoords[0], y: oldCoords[1]}].pheromone += 1
						}
					}
					updates++
				}
			}

			if updates > 0 {
				roundCounter++
			} else {
				roundCounter = 0

				for ant := range ams.Routes {
					coords := ams.IDtoMap[ams.Routes[ant][iterationCounter][len(ams.Routes[ant][iterationCounter])-1]]
					ams.cells[coordinate{x: coords[0], y: coords[1]}].ants -= 1
				}

				if iterationCounter >= len(ams.Routes[0]) {
					fmt.Println()
					over = true
				} else {
					iterationCounter++
					pheroActivate = make([]bool, len(ams.Routes))
					ams.deprecatePheromones()
				}
			}
		}

	}
}

// Called when system is initialized
func (ams *AntMoverSystem) New(world *ecs.World) {
	resetAll()
	ams.cells = make(map[coordinate]cellEntity)
	pheroActivate = make([]bool, len(ams.Routes))
	roundTimer = ams.SecondsPerRound
}

// Called to add entities to the system
func (ams *AntMoverSystem) Add(basic *ecs.BasicEntity, render *common.RenderComponent, cell *CellComponent, space *common.SpaceComponent, x int, y int) {
	ams.cells[coordinate{x: x, y: y}] = cellEntity{*basic, render, space, cell}
}

func updateRoundTimer(dt float32) {
	roundTimer += dt
}

func resetAll() {
	roundTimer = 0
	roundCounter = 0
	iterationCounter = 0
	over = false
}

func (ams *AntMoverSystem) deprecatePheromones() {
	for _, cell := range ams.cells {
		if cell.pheromone > 0 {
			cell.pheromone *= 0.5
		}
	}
}
