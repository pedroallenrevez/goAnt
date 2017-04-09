package world

import (
	"github.com/op/go-logging"
	"strconv"
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
	x          int
	y          int
}

// Necessary because of the way pointers to structs in maps work
// see: http://stackoverflow.com/q/32751537
func (c cell) incrementAnts() {
	c.ants++
}

// See incrementAnts
func (c cell) decrementAnts() {
	c.ants--
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
	Init(int, func(Point, Point) float64)
}

// WorldImpl Implementation of the interface World
type WorldImpl struct {
	//See what's needed
	antMap          map[NodeID]cell
	pheroMap        map[nodePair]PheromoneValue
	updatedPheroMap map[nodePair]PheromoneValue
	distance        calculateDistance
}

//Init Initializes world with matrix size and a distance function
func (w *WorldImpl) Init(size int, funcArg func(Point, Point) float64) {
	//generate antMap
	antMap := make(map[NodeID]cell, size*size)

	concatInt := func(x, y int) NodeID {
		i, err := strconv.Atoi(strconv.Itoa(x) + strconv.Itoa(y))
		if err != nil {
			panic("could not convert indexes")
		}
		return NodeID(i)
	}
	calcNeighbours := func(x, y, size int) []NodeID {
		slice := make([]NodeID, 8)
		if x-1 >= 0 {
			// -1 -1  -1 0  -1 1
			slice = append(slice, concatInt(x-1, y))
			if y-1 >= 0 {
				slice = append(slice, concatInt(x-1, y-1))
				slice = append(slice, concatInt(x, y-1))

			}
			if y+1 <= size-1 {
				slice = append(slice, concatInt(x-1, y+1))
				slice = append(slice, concatInt(x, y+1))
			}

		}
		if x+1 <= size-1 {
			slice = append(slice, concatInt(x+1, y))
			if y-1 >= 0 {
				slice = append(slice, concatInt(x+1, y-1))
			}
			if y+1 <= size-1 {
				slice = append(slice, concatInt(x+1, y+1))
			}

		}
		return slice

	}

	for x := 0; x < size-1; x++ {
		for y := 0; y < size-1; y++ {
			index := concatInt(x, y)
			newCell := cell{
				id:         index,
				ants:       0,
				x:          x,
				y:          y,
				neighbours: calcNeighbours(x, y, size),
			}
			antMap[index] = newCell

		}
	}
	w.antMap = antMap
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
	log.Error("not yet implemented")
	return false

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

func (w *WorldImpl) getCell(node NodeID) cell {
	if current, ok := w.antMap[node]; ok {
		return current
	}

	log.Critical("Nodeid %d does not exist! Of course i can't get this cell!", node)
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

func (w *WorldImpl) computeDistance(start, end cell) float64 {
	p1 := Point{start.x, start.y}
	p2 := Point{end.x, end.y}
	return w.distance(p1, p2)
}
