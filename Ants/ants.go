package Ant

import (
	"fmt"
	"github.com/pedroallenrevez/world"
	"math"
	"math/rand"
	"time"
)

//Ant badjriziz
type Ant struct {
	location      NodeID  // present location of ant
	route         []int   // list of navigated nodes
	alpha         float64 // experimental value - influences pheromones
	beta          float64 // experimental value - influences desirability
	distTravelled float64 // total distance travelled
	tourComplete  bool    // ? ant reached goal
	firstPass     bool    // ? first iteration
	world         World   // * -> worldMap
}

//initializes ants and resets them for new iteration
func (ant *Ant) init() {
	ant.location = 0
	ant.route = nil
	append(ant.route, 0)
	ant.alpha = 1.0
	ant.distTravelled = 0
	tourComplete = false
	//set world
}

// pickPath
// randomly selects a path based on the ACO algorithm
func (ant *Ant) pickPath() Node {
	// termination cases
	// if node not connected to anyone but last
	var path int
	var nodes = world.possibleMoves(ant.location)
	if ant.firstPass {
		//choose randomly even from all nodes
		return nodes[randInt(0, len(nodes))]
	}

	interval := make([]float64)
	accumulate = 0
	for i, node := range nodes {
		probability := math.Pow(node.pheromone, ant.alpha) * //priori
			math.Pow(1.0/node.distance, ant.beta) //posterior
		append(interval, probability)
		accumulate += probability
	}
	random = randFloat() * accumulate
	for i := range interval {
		nodes[i] /= accumulate
		if random < interval[i] {
			//i dont have the id!
			//dont need cus indexes match
			return nodes[i]
		}
	}

}

// traverse
// updates position on map and ant, updates route, update distance
func (ant *Ant) traverse(start, end int) {
	world.updatePosition(start, end)
	ant.location = end
	//route ++ end
	append(ant.route, end)
	//dist += dist callback
	ant.distTravelled += world.calculateDistance(start, end)

}

//do I need this? just eliminate loops on world
func (ant *Ant) isFood(location int) bool {
	//reached food
	//switch mode
	//ant.forwardMode = false
	//delete loops
	if world.isFood(ant.location) {
		//delete loops
		/*for i := range ant.route {
			// ? cross correct loops and distance travelled
			//floyd loop deletion
		}*/
	}
}

// returns a list of node ids
func (ant *Ant) run() []int {
	//GO ROUTINE
	//run until tour is complete
	for !tourComplete {
		//next node
		var next = pickPath()
		traverse(ant.location, next.id)
		if isFood(next) {
			//tour complete
			tourComplete = true
		}
	}
	return ant.route
}

func randInt(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func randFloat() float64 {
	rand.Seed(time.Now().Unix())
	return rand.Float64()
}
