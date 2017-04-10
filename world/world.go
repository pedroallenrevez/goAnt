package world

import (
	"github.com/op/go-logging"
	"math/rand"
)

var log = logging.MustGetLogger("example")

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

const initialPheromoneValue = 1
const decayFactor = PheromoneValue(0.5)

// NodeID abstraction for index
type NodeID int

// PheromoneValue abstraction for float64
type PheromoneValue float64
type calculateDistance func(Point, Point) float64

// Point abstraction for calculating distances
type Point struct {
	X, Y int
}

// Node is returned to the ant when possibleMoves is called
type Node struct {
	Pheromone PheromoneValue
	Distance  float64
	ID        NodeID
}

// Cell is the internal type for the representation of graph nodes
type cell struct {
	id         NodeID
	ants       int
	neighbours []NodeID
	goal       bool
	x          int
	y          int
}

// Necessary because of the way pointers to structs in maps work
// see: http://stackoverflow.com/q/32751537
func (c *cell) incrementAnts() {
	c.ants++
}

// See incrementAnts
func (c *cell) decrementAnts() {
	c.ants--
}

func (c *cell) setGoal() {
	c.goal = true
}

// nodePair is an abstraction to be used as a key in the pheromone maps
type nodePair struct {
	previous NodeID
	next     NodeID
}

func (p nodePair) invert() nodePair {
	return nodePair{p.next, p.previous}
}

// World defines the interface with which the ants interact
type World interface {
	PossibleMoves(NodeID) []Node
	UpdatePosition(NodeID, NodeID)
	PutPheromone(NodeID, NodeID)
	UpdatePheromones()
	IsGoal(NodeID) bool
	ExportMap() [][]int
	Init(int, func(Point, Point) float64)
}

// WorldImpl Implementation of the interface World
type WorldImpl struct {
	//See what's needed
	antMap          map[NodeID]*cell
	pheroMap        map[nodePair]PheromoneValue
	updatedPheroMap map[nodePair]PheromoneValue
	distance        calculateDistance
}

//Init Initializes world with matrix size and a distance function
func (w *WorldImpl) Init(size int, funcArg func(Point, Point) float64) {
	//generate aNtMap
	antMap := make(map[NodeID]*cell)

	//starts at -1 so ID 0 is the first
	startingID := -1

	createUniqueID := func() NodeID {
		startingID++
		return NodeID(startingID)
	}

	calcNeighbours := func(x, y, size int, node int) []NodeID {
		slice := make([]NodeID, 0, 8)
		if x-1 >= 0 {
			//This math is boring, trust me, as long as they are generated
			//1 line at a time, this is correct ;)
			slice = append(slice, NodeID(node-1))
			if y-1 >= 0 {
				slice = append(slice, NodeID(node-size-1))

			}
			if y+1 <= size-1 {
				slice = append(slice, NodeID(node+size-1))
			}

		}
		if y-1 >= 0 {
			slice = append(slice, NodeID(node-size))
		}
		if y+1 <= size-1 {
			slice = append(slice, NodeID(node+size))
		}
		if x+1 <= size-1 {
			slice = append(slice, NodeID(node+1))
			if y-1 >= 0 {
				slice = append(slice, NodeID(node-size+1))
			}
			if y+1 <= size-1 {
				slice = append(slice, NodeID(node+size+1))
			}

		}
		return slice

	}

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			index := createUniqueID()
			antMap[index] = &cell{
				id:         index,
				ants:       0,
				x:          x,
				y:          y,
				neighbours: calcNeighbours(x, y, size, int(index)),
			}

		}
	}

	//set one cell as goal
	goal := NodeID(rand.Intn(startingID - 1))
	antMap[goal].setGoal()
	log.Warning("Goal is ", goal)
	log.Warning(antMap[goal].goal)

	w.antMap = antMap

	for _, cell := range w.antMap {
		log.Debug(*cell)
	}

	for _, cell := range w.ExportMap() {
		log.Debug(cell)
	}
	//generate according size pheroMap map nodepair -> pheromone
	w.pheroMap = make(map[nodePair]PheromoneValue)
	w.updatedPheroMap = make(map[nodePair]PheromoneValue)
	w.distance = funcArg
	//copy one to updatedPheroMap
	//assign function
}

// PossibleMoves Given a NodeID returns the possible moves the ant can make from
// that position, this includes pheromonal, distance and ID values
// for every possible movement.
func (w *WorldImpl) PossibleMoves(nodeid NodeID) []Node {
	current := w.getCell(nodeid)

	var result = make([]Node, len(current.neighbours))

	for i, neighbour := range current.neighbours {
		nCell := w.getCell(neighbour)
		result[i].ID = nCell.id
		result[i].Distance = w.computeDistance(current, nCell)
		result[i].Pheromone = w.getPheromone(current.id, nCell.id)
	}

	return result
}

// UpdatePosition Used by the ants to update their position on map
func (w *WorldImpl) UpdatePosition(before NodeID, after NodeID) {
	w.getCell(before).incrementAnts()
	w.getCell(after).decrementAnts()
}

// PutPheromone deprecated
// not needed for ant
// ant will return route on run. just process route and deposit pheromone then
// can be used by real time ants
func (w *WorldImpl) PutPheromone(before NodeID, after NodeID) {
	pair := nodePair{before, after}
	w.updatedPheroMap[pair] += initialPheromoneValue
}

// IsGoal Returns true if the given node is a goal node
func (w *WorldImpl) IsGoal(node NodeID) bool {
	return w.getCell(node).goal

}

// UpdatePheromones decays the map with decay constant
func (w *WorldImpl) UpdatePheromones() {
	for pair, pheromone := range w.pheroMap {
		pheromone *= decayFactor
		if updatedVal, ok := w.updatedPheroMap[pair]; ok {
			pheromone += updatedVal
		} else if updatedVal, ok := w.updatedPheroMap[pair.invert()]; ok {
			pheromone += updatedVal
		}
	}
}

func (w *WorldImpl) getCell(node NodeID) *cell {
	if current, ok := w.antMap[node]; ok {
		return current
	}

	log.Critical("Nodeid", node, " does not exist! Of course i can't get this cell!")
	//TODO
	panic("This should never have happened... CALL A MEDIC!")
}

func (w *WorldImpl) getPheromone(start NodeID, end NodeID) PheromoneValue {
	pair := nodePair{start, end}
	result := PheromoneValue(0.0)

	if pheromone, ok := w.pheroMap[pair]; ok {
		result = pheromone
	} else if pheromone, ok := w.pheroMap[pair.invert()]; ok {
		result = pheromone
	}
	return result
}

func (w *WorldImpl) computeDistance(start, end *cell) float64 {
	p1 := Point{start.x, start.y}
	p2 := Point{end.x, end.y}
	return w.distance(p1, p2)
}

func (w WorldImpl) ExportMap() [][]int {
	if w.antMap == nil {
		panic("No map yet, CANT EXPORT!")
	}

	size := 0

	for _, cell := range w.antMap {
		if cell.x > size {
			size = cell.x
		}
		if cell.y > size {
			size = cell.y
		}
	}

	log.Debug("AntMap has size:", size)

	result := make([][]int, size+1)

	for i := range result {
		result[i] = make([]int, size+1)
	}

	for _, cell := range w.antMap {
		// if goal = 1, normal = 0, obstacle = 2
		if cell.goal {
			result[cell.x][cell.y] = 1
		} else {
			result[cell.x][cell.y] = 0
		}
	}

	return result

}
