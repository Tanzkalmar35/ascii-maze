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
	xDistance := float64(node.cell.x - n.cell.x)
	yDistance := float64(node.cell.y - n.cell.y)
	return int(math.Abs(xDistance + yDistance))
}

type SolvingStrategy struct {
}

func (ss SolvingStrategy) AStar(m Maze) []Cell {
	var openNodes []Node
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
		time.Sleep(20 * time.Millisecond)
	}
}

func (ss SolvingStrategy) DeadEndFilling(m Maze) {
	var closedCells []Cell
	for y, _ := range m.grid {
		for x, _ := range m.grid[y] {
			if m.grid[y][x] != Path {
				continue
			}

			var changedCells []Cell
			changedCells = append(changedCells, Cell{x: x, y: y})
			neighbours := NewDefaultNode(changedCells[0]).GetNeighbours(m)
			numberOfClosedNodesSurrounding := 0

			for _, neighbour := range neighbours {
				if m.grid[neighbour.cell.y][neighbour.cell.x] == Wall {
					numberOfClosedNodesSurrounding++
				} else if Contains(closedCells, neighbour.cell) {
					numberOfClosedNodesSurrounding++
				}
			}
			if numberOfClosedNodesSurrounding >= 3 {
				closedCells = append(closedCells, changedCells[0])
				// Paint that cell
				m.grid[changedCells[0].y][changedCells[0].x] = Indicator
				m.UpdateDisplay(changedCells)
			}
		}
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
