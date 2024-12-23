package main

import (
	"math"
	"slices"
	"time"
)

var endNode Node

type Node struct {
	cell   Cell
	fCost  int
	hCost  int
	gCost  int
	parent *Node
}

func NewDefaultNode(cell Cell) Node {
	return Node{cell: cell}
}

func NewNode(cell Cell, hCost, gCost int, parent *Node) Node {
	fCost := hCost + gCost
	return Node{cell: cell, fCost: fCost, hCost: hCost, gCost: gCost, parent: parent}
}

func EndNode() Node {
	return endNode
}

func (n Node) GetNeighbours(m Maze) []Node {
	var neighbours []Node
	if m.IsPathInMaze(Cell{x: n.cell.x - 1, y: n.cell.y}) {
		cell := Cell{x: n.cell.x - 1, y: n.cell.y}
		newNode := NewNode(cell, cell.GetDistanceToNode(EndNode()), n.gCost+1, &n)
		neighbours = append(neighbours, newNode)
	}
	if m.IsPathInMaze(Cell{x: n.cell.x + 1, y: n.cell.y}) {
		cell := Cell{x: n.cell.x + 1, y: n.cell.y}
		newNode := NewNode(cell, cell.GetDistanceToNode(EndNode()), n.gCost+1, &n)
		neighbours = append(neighbours, newNode)
	}
	if m.IsPathInMaze(Cell{x: n.cell.x, y: n.cell.y - 1}) {
		cell := Cell{x: n.cell.x, y: n.cell.y - 1}
		newNode := NewNode(cell, cell.GetDistanceToNode(EndNode()), n.gCost+1, &n)
		neighbours = append(neighbours, newNode)
	}
	if m.IsPathInMaze(Cell{x: n.cell.x, y: n.cell.y + 1}) {
		cell := Cell{x: n.cell.x, y: n.cell.y + 1}
		newNode := NewNode(cell, cell.GetDistanceToNode(EndNode()), n.gCost+1, &n)
		neighbours = append(neighbours, newNode)
	}
	return neighbours
}

func (n Node) GetDistance(node Node) int {
	xDistance := float64(n.cell.x - node.cell.x)
	yDistance := float64(n.cell.y - node.cell.y)
	return int(math.Abs(xDistance + yDistance))
}

type SolvingStrategy struct {
}

func (ss SolvingStrategy) AStarAlgo(m Maze) []Cell {
	var openNodes []Node
	var closedNodes []Node
	var currentNode Node

	endNode = NewDefaultNode(Cell{x: m.end.x, y: m.end.y})
	startingNode := NewNode(m.start, m.start.GetDistanceToNode(EndNode()), 0, nil)

	openNodes = append(openNodes, startingNode)
	currentNode = startingNode

	for {
		if len(openNodes) == 0 {
			panic("WIOUDALKSJDLKA")
		}

		currentNode = openNodes[0]

		// select the next node
		for _, node := range openNodes {
			if node.fCost < currentNode.fCost {
				currentNode = node
			} else if node.fCost == currentNode.fCost &&
				node.hCost < currentNode.hCost {
				currentNode = node
			}
		}
		closedNodes = append(closedNodes, currentNode)
		openNodes = Remove(openNodes, currentNode)

		if currentNode.GetDistance(EndNode()) == 0 {
			currentPathTile := EndNode()
			var path []Cell
			for {
				if currentNode.GetDistance(startingNode) == 0 {
					break
				}

				path = append(path, m.start)
				path = append(path, m.end)
				m.grid[m.start.y][m.start.x] = Indicator
				m.grid[m.end.y][m.end.x] = Indicator
				currentPathTile = *currentPathTile.parent
			}

			m.UpdateDisplay(path)
			return path
		}

		// Get neighbours of the current node
		neighbours := currentNode.GetNeighbours(m)

		for _, neighbour := range neighbours {
			costOfNavigatingToNeighbour := currentNode.gCost + currentNode.GetDistance(neighbour)
			neighbourIsOpen := Contains(openNodes, neighbour)
			if !neighbourIsOpen || costOfNavigatingToNeighbour < neighbour.gCost {
				neighbour.gCost = costOfNavigatingToNeighbour
				neighbour.parent = &currentNode

				if !neighbourIsOpen {
					neighbour.hCost = neighbour.GetDistance(EndNode())
					openNodes = append(openNodes, neighbour)
				}
			}
		}

		var path []Cell

		for _, node := range closedNodes {
			m.grid[node.cell.y][node.cell.x] = Indicator
			path = append(path, node.cell)
		}

		m.UpdateDisplay(path)

		time.Sleep(20 * time.Millisecond)
	}

	return nil
}

func RetracePath(startingNode, targetNode Node, m Maze) {
	var path []Node
	var changedCells []Cell

	currentNode := targetNode

	for {
		if currentNode != startingNode {
			path = append(path, currentNode)
			m.grid[currentNode.cell.y][currentNode.cell.x] = Indicator
			changedCells = append(changedCells, currentNode.cell)
			currentNode = *currentNode.parent
		} else {
			m.UpdateDisplay(changedCells)
			break
		}
	}

	slices.Reverse(path)
}
