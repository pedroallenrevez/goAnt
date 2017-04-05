package Ant

import (
	"fmt"
)

//Ant badjriziz
type Ant struct {
	location      int
	route         []int
	alpha         float64
	distTravelled float64
	forwardMode   bool
	//working mode
}

func (ant *Ant) pickPath() {
	//OR
	switch ant.forwardMode {
	case True:
		//probability
	case False:
		//retrace route
		//drop pheromone
	}
}

func (ant *Ant) traverse(start, end int) {
	ant.location = end
	//route ++ end
	append(ant.route, end)
	//dist += dist callback

}

func (ant *Ant) isFood(location int) bool {
	//switch mode
	//delete loops
}

func (ant *Ant) isNest(location int) bool {
	//switch mode
	//delete loops
}

func (ant *Ant) run(possMoves []int) {
	//run until tour is complete

	//while possible moves
	//pickpath
	//traverse
	var next = pickPath()
	traverse(ant.location, next)
	if isFood(next) {
		//switch modes
	}
	if isNest(next) {
		//finish tour
	}
}
