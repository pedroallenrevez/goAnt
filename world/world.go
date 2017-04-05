package main

import (
	"fmt"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("example")

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

type NodeID int
type PheromoneValue float64
type calculateDistance func(cell, cell) float64

// Node is returned to the ant when possibleMoves is called
type Node struct {
	pheromone PheromoneValue
	distance  float64
	id        int
}

// Cell is the internal type for the representation of graph nodes
type cell struct {
	ants       int
	neighbours []NodeID
	x          int
	y          int
}

// nodePair is an abstraction to be used as a key in the pheromone maps
type nodePair struct {
	previous NodeID
	next     NodeID
}

// World defines the interface with which the ants interact
type World interface {
	possibleMoves(NodeID) []Node
	updatePosition(NodeID, NodeID)
	putPheromone(NodeID, NodeID)
}

// Implementation of the interface World
type WorldImpl struct {
	//See what's needed
	antMap          map[NodeID]cell
	pheroMap        map[nodePair]PheromoneValue
	updatedPheroMap map[nodePair]PheromoneValue
	distance        calculateDistance
}

// Given a NodeID returns the possible moves the ant can make from
// that position, this includes pheromonal, distance and ID values
// for every possible movement.
func (w worldImpl) possibleMoves(nodeid NodeID) []Node {
	current = getCell(nodeid)

	var result [len(current)]Node

	for i, neighbour := range current.neighbours {
		result[i].distance = w.distance(current, getCell(neighbour))
		result[i].pheromone = w.
	}
}

func (w worldImpl) updatePosition(before NodeID, after NodeID) {

}

func (w worldImpl) putPheromone(before NodeID, after NodeID) {

}

func (w worldImpl) getCell(node NodeID) cell {
	// Trick to get zero value of cell, surelly the must be another way to do this
	cellZero := cell

	current := w.antMap[node]

	if current != cellzero {
		return current
	}

	log.Critical("Nodeid %d does not exist! Of course i can't get this cell!", nodeid)
	return nil
}

func (w worldImpl) getPheromone (start NodeID, end NodeID) PheromoneValue {
	// Trick to get zero value of , surelly the must be another way to do this
	cellZero := cell

	current := w.antMap[node]

	if current != cellzero {
		return current
	}

	log.Critical("Nodeid %d does not exist! Of course i can't get this cell!", nodeid)
	return nil
}
