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

const initialPheromoneValue = 1
const decayFactor = PheromoneValue(0.5)

// NodeID abstraction for index
type NodeID int

// PheromoneValue abstraction for float64
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
	possibleMoves(NodeID) []Node
	updatePosition(NodeID, NodeID)
	putPheromone(NodeID, NodeID)
}

// WorldImpl Implementation of the interface World
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
	w.getCell(before).incrementAnts()
	w.getCell(after).decrementAnts()
}

func (w WorldImpl) putPheromone(before NodeID, after NodeID) {
	pair := nodePair{before, after}
	w.updatedPheroMap[pair] += initialPheromoneValue
}

func (w WorldImpl) updatePheromones() {
	for pair, pheromone := range w.pheroMap {
		pheromone *= decayFactor
		if updatedVal, ok := w.updatedPheroMap[pair]; ok {
			pheromone += updatedVal
		} else if updatedVal, ok := w.updatedPheroMap[pair.invert()]; ok {
			pheromone += updatedVal
		}
	}
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

// Create Create a blank world and return the addr
// Factory to produce interfaces for the ants to use
func (w worldImpl) Create() *World {
	return &worldImpl{
		antMap:          w.antMap,
		pheroMap:        w.pheroMap,
		updatedPheroMap: w.updatedPheroMap,
	}
}
