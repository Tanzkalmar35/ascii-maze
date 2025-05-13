package main

import (
	"math"
	"slices"
	"time"
)

var solving_strategy string
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
	xDistance := float64(node.cell.x - n.cell.x)
	yDistance := float64(node.cell.y - n.cell.y)
	return int(math.Abs(xDistance + yDistance))
}

type SolvingStrategy struct {
}

func (ss SolvingStrategy) AStar(m Maze) []Cell {
	var openNodes []Node
	m.SetStrategy("astar")
	endNode = NewDefaultNode(Cell{x: m.end.x, y: m.end.y})
	startingNode := NewNode(m.start, m.start.GetDistanceToNode(EndNode()), 0, nil)

	openNodes = append(openNodes, startingNode)

	for {
		if len(openNodes) == 0 {
			panic("No path found!")
		}

		// Select the node with the lowest fCost
		minIndex := 0
		for i, node := range openNodes {
			if node.fCost < openNodes[minIndex].fCost ||
				(node.fCost == openNodes[minIndex].fCost && node.hCost < openNodes[minIndex].hCost) {
				minIndex = i
			}
		}

		currentNode := openNodes[minIndex]
		openNodes = append(openNodes[:minIndex], openNodes[minIndex+1:]...) // Remove the selected node

		// Check if we reached the target
		if currentNode.cell.x == EndNode().cell.x && currentNode.cell.y == EndNode().cell.y ||
			currentNode.cell.x == EndNode().cell.x+1 && currentNode.cell.y == EndNode().cell.y ||
			currentNode.cell.x == EndNode().cell.x-1 && currentNode.cell.y == EndNode().cell.y ||
			currentNode.cell.x == EndNode().cell.x && currentNode.cell.y == EndNode().cell.y-1 ||
			currentNode.cell.x == EndNode().cell.x && currentNode.cell.y == EndNode().cell.y+1 {
			return RetracePath(startingNode, currentNode, m) // Return the path
		}

		// Get neighbours of the current node
		neighbours := currentNode.GetNeighbours(m)

		for _, neighbour := range neighbours {
			costOfNavigatingToNeighbour := currentNode.gCost + currentNode.GetDistance(neighbour)
			neighbourIsOpen := Contains(openNodes, neighbour)

			if neighbourIsOpen {
				// Update the existing node in openNodes
				for i, openNode := range openNodes {
					if openNode.cell == neighbour.cell {
						if costOfNavigatingToNeighbour < openNode.gCost {
							openNodes[i].gCost = costOfNavigatingToNeighbour
							openNodes[i].parent = &currentNode
						}
						break
					}
				}
			} else {
				// Add the new neighbour to openNodes
				neighbour.hCost = neighbour.GetDistance(EndNode())
				openNodes = append(openNodes, neighbour)
			}
		}

		// Optionally update the display for debugging
		var path []Cell
		for _, node := range openNodes {
			m.grid[node.cell.y][node.cell.x] = Indicator
			path = append(path, node.cell)
		}
		m.UpdateDisplay(path)
		time.Sleep(1 * time.Millisecond)
	}
}

func (ss SolvingStrategy) DeadEndFilling(m Maze) {
	var closedCells []Cell
	m.SetStrategy("deadendfilling")

	for {
		var changedCells []Cell
		var tempClosedCells []Cell
		anyChange := false // Flag to track if any changes were made in this iteration

		for i := 0; i < len(m.grid); i++ { // Loop through rows
			for j := 0; j < len(m.grid[i]); j++ {
				if m.grid[i][j] != Path {
					continue
				}

				changedCells = append(changedCells, Cell{x: j, y: i})
				neighbours := NewDefaultNode(changedCells[0]).GetNeighbours(m)
				amountOfNeighbours := len(neighbours)

				// Count only the neighbors that are in closedCells
				for _, neighbour := range neighbours {
					if Contains(closedCells, neighbour.cell) {
						amountOfNeighbours--
					}
				}

				// If it's a dead end, mark it
				if amountOfNeighbours <= 1 {
					tempClosedCells = append(tempClosedCells, changedCells[0])
					anyChange = true // Indicate that a change was made
				}

				changedCells = changedCells[:0] // Clear changedCells for the next iteration
			}
		}

		// Paint the cells
		for _, cell := range tempClosedCells {
			m.grid[cell.y][cell.x] = Indicator
		}

		// Update closedCells with the newly closed cells
		closedCells = append(closedCells, tempClosedCells...)

		// If no changes were made, break the loop
		if !anyChange {
			break
		}

		// Optionally update the display here if needed
		m.UpdateDisplay(closedCells)
	}
}

func RetracePath(startingNode, targetNode Node, m Maze) []Cell {
	var path []Cell
	currentNode := targetNode

	for currentNode != startingNode {
		path = append(path, currentNode.cell)
		m.grid[currentNode.cell.y][currentNode.cell.x] = BestPath

		if currentNode.parent == nil {
			break // Prevent dereferencing nil
		}
		currentNode = *currentNode.parent
	}

	// Reverse the path to get it from start to end
	slices.Reverse(path)
	return path
}
