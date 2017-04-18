package ants

import (
	"fmt"
	"github.com/pedroallenrevez/goAnt/world"
	"math"
	"math/rand"
	"time"
)

//AntInterface Specifies interactability for Ant
type AntInterface interface {
	Reset()
	Run(chan []world.NodeID)
}

//Ant badjriziz
type Ant struct {
	location      world.NodeID // present location of ant
	previous      world.NodeID
	route         []world.NodeID // list of navigated nodes
	alpha         float64        // experimental value - influences pheromones
	beta          float64        // experimental value - influences desirability
	distTravelled float64        // total distance travelled
	tourComplete  bool           // ? ant reached goal
	firstPass     bool           // ? first iteration
	world         world.World    // * -> worldMap
}

//Reset reset ant for new iteration
func (ant *Ant) Reset() {
	ant.location = 0
	ant.route = nil
	ant.alpha = 0.0
	ant.distTravelled = 0
	ant.tourComplete = false
	ant.firstPass = true
}

//NewAnt constructor for a new ant
func NewAnt(w world.World) *Ant {
	new := Ant{
		alpha:        1.0,
		firstPass:    true,
		tourComplete: false,
		world:        w,
	}
	//new.route = append(new.route, 0)

	return &new
}

// pickPath
// randomly selects a path based on the ACO algorithm
func (ant *Ant) pickPath() world.Node {
	// termination cases
	// if node not connected to anyone but last
	var nodes = ant.world.PossibleMoves(ant.location)
	fmt.Println("possibles moves", nodes)
	var previous world.Node
	if nodes != nil {
		for i, node := range nodes {
			if node.ID == ant.previous {
				//slice trick to delete element at i position
				previous = node
				nodes = append(nodes[:i], nodes[i+1:]...)
			}
		}
	}

	interval := make([]float64, len(nodes))
	accumulate := 0.0
	for i, node := range nodes {
		probability := math.Pow(float64(node.Pheromone), ant.alpha) * //priori
			math.Pow(1.0/node.Distance, ant.beta) //posterior
		accumulate += probability
		interval[i] = accumulate
	}

	random := randFloat()
	for i := range interval {
		interval[i] /= accumulate
		//fmt.Println(interval[i])
		//fmt.Println(random)
		if random < interval[i] {
			//fmt.Println("chosen node")
			//fmt.Println(nodes[i])
			return nodes[i]
		}
	}

	fmt.Println("REACHING THE ENDREACHING THE END")
	return previous
}

// traverse
// updates position on map and ant, updates route, update distance
func (ant *Ant) traverse(start world.NodeID, end world.Node) {
	ant.previous = ant.location
	ant.location = end.ID
	//route ++ end
	ant.route = append(ant.route, end.ID)
	//dist += dist callback
	ant.distTravelled += end.Distance
	ant.world.UpdatePosition(ant.location, end.ID)

}

// if this node is a goal complete route
func (ant *Ant) isGoal(node world.NodeID) {

	if ant.world.IsGoal(node) {
		ant.tourComplete = true
	}
}

//Run returns a list of node ids
func (ant *Ant) Run(channel chan []world.NodeID) {
	//run until tour is complete
	for !ant.tourComplete {
		var next = ant.pickPath()
		ant.traverse(ant.location, next)
		ant.isGoal(next.ID)
	}
	//loop deletion
	ant.route = ant.loopDeletion(ant.route)
	ant.firstPass = false

	channel <- ant.route
}

func (ant *Ant) loopDeletion(nodes []world.NodeID) []world.NodeID {
	result := make([]world.NodeID, 0)
	for _, node := range nodes {
		if index, ok := inList(node, result); ok {
			//delete
			result = result[:index+1]
		} else {
			//add
			result = append(result, node)
		}
	}
	return result
}

func randInt(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func randFloat() float64 {
	rand.Seed(time.Now().Unix())
	return rand.Float64()
}

func inList(val world.NodeID, arr []world.NodeID) (int, bool) {
	for i, node := range arr {
		if val == node {
			return i, true
		}
	}
	return 0, false
}
