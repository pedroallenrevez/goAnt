package ants

import (
	//"fmt"
	"github.com/pedroallenrevez/goAnt/world"
	"math"
	"math/rand"
	"time"
)

//Ant badjriziz
type Ant struct {
	location      world.NodeID // present location of ant
	route         []int        // list of navigated nodes
	alpha         float64      // experimental value - influences pheromones
	beta          float64      // experimental value - influences desirability
	distTravelled float64      // total distance travelled
	tourComplete  bool         // ? ant reached goal
	firstPass     bool         // ? first iteration
	world         world.World  // * -> worldMap
}

//initializes ants and resets them for new iteration
func (ant *Ant) init(w world.World) {
	ant.location = 0
	ant.route = nil
	ant.route = append(ant.route, 0)
	ant.alpha = 1.0
	ant.distTravelled = 0
	ant.tourComplete = false
	ant.world = w
	//set world
}

// pickPath
// randomly selects a path based on the ACO algorithm
func (ant *Ant) pickPath() world.Node {
	// termination cases
	// if node not connected to anyone but last
	var nodes = ant.world.PossibleMoves(ant.location)
	if ant.firstPass {
		//choose randomly even from all nodes
		return nodes[randInt(0, len(nodes))]
	}

	interval := make([]float64, len(nodes))
	accumulate := 0.0
	for _, node := range nodes {
		probability := math.Pow(float64(node.Pheromone), ant.alpha) * //priori
			math.Pow(1.0/node.Distance, ant.beta) //posterior
		interval = append(interval, probability)
		accumulate += probability
	}
	random := randFloat() * accumulate
	for i := range interval {
		interval[i] /= accumulate
		if random < interval[i] {
			//i dont have the id!
			//dont need cus indexes match
			return nodes[i]
		}
	}
	//FIX
	//thors hammers
	return nodes[0]
}

// traverse
// updates position on map and ant, updates route, update distance
func (ant *Ant) traverse(start world.NodeID, end world.Node) {
	ant.location = end.ID
	//route ++ end
	ant.route = append(ant.route, int(end.ID))
	//dist += dist callback
	ant.distTravelled += end.Distance
	ant.world.UpdatePosition(ant.location, end.ID)

}

//TODO - is goal function!
func (ant *Ant) isGoal(node world.NodeID) {
	if ant.world.IsGoal(node) {
		ant.tourComplete = true
	}
	//channel trigger
}

// returns a list of node ids
func (ant *Ant) run() []int {
	//GO ROUTINE
	//run until tour is complete
	for !ant.tourComplete {
		//next node
		var next = ant.pickPath()
		ant.traverse(ant.location, next)
		ant.isGoal(next.ID)
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
