package main

import (
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
	id        NodeID
}

// Cell is the internal type for the representation of graph nodes
type cell struct {
	id         NodeID
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

func (p nodePair) invert() nodePair {
	return nodePair{p.next, p.previous}
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
func (w WorldImpl) possibleMoves(nodeid NodeID) []Node {
	current := w.getCell(nodeid)

	var result = make([]Node, len(current.neighbours))

	for i, neighbour := range current.neighbours {
		nCell := w.getCell(neighbour)
		result[i].id = nCell.id
		result[i].distance = w.distance(current, nCell)
		result[i].pheromone = w.getPheromone(current.id, nCell.id)
	}

	return result
}

func (w WorldImpl) updatePosition(before NodeID, after NodeID) {

}

func (w WorldImpl) putPheromone(before NodeID, after NodeID) {

}

func (w WorldImpl) getCell(node NodeID) cell {
	if current, ok := w.antMap[node]; ok {
		return current
	}

	log.Critical("Nodeid %d does not exist! Of course i can't get this cell!", node)
	//TODO
	panic("This should never have happened... CALL A MEDIC!")
}

func (w WorldImpl) getPheromone(start NodeID, end NodeID) PheromoneValue {
	pair := nodePair{start, end}
	result := PheromoneValue(0.0)

	if pheromone, ok := w.pheroMap[pair]; ok {
		result = pheromone
	} else if pheromone, ok := w.pheroMap[pair.invert()]; ok {
		result = pheromone
	}
	return result
}
