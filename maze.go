package main

import (
	"math/rand"
	"time"
)

const (
	Wall      = '█'
	Path      = ' '
	Start     = 'S'
	End       = 'E'
	Indicator = '*'
)

var currentMaze *Maze

type Cell struct {
	x, y int
}

type Maze struct {
	height, width int
	grid          [][]rune // Beware, its yx not xy
}

func DisplayMaze(width, height int) *Maze {
	maze := &Maze{width: width, height: height}
	maze.grid = make([][]rune, height)
	// Loop vertically
	for i := range maze.grid {
		maze.grid[i] = make([]rune, width)
		// Loop horizontally
		for j := range maze.grid[i] {
			maze.grid[i][j] = Wall
		}
	}
	return maze
}

func (m *Maze) Render() string {
	var view string
	for i := range m.grid {
		for j := range m.grid[i] {
			switch m.grid[i][j] {
			case Wall:
				view += "\033[38;5;15m" + string(Wall) + "\033[0m" // White color for walls
			case Path:
				view += "\033[38;5;15m" + string(Path) + "\033[0m" // White color for paths
			case Start:
				view += "\033[38;5;15m" + string(Start) + "\033[0m" // White color for start
			case End:
				view += "\033[38;5;15m" + string(End) + "\033[0m" // White color for end
			case Indicator:
				view += "\033[38;5;14m" + string(Indicator) + "\033[0m" // Yellow color for indicator
			}
		}
		view += "\n"
	}
	return view
}

func (m *Maze) GenerateRandomPrimMaze() {
	var frontier []Cell

	// Generate random starting point inside the grid if its the first iteration
	frontier = m.GenerateFirstPrimIteration()

	// Paint frontiers
	for _, f := range frontier {
		m.grid[f.y][f.x] = Indicator
	}

	for {
		if len(frontier) <= 0 {
			break
		}

		// Choose one of the frontiers randomly
		randomlyChosenFrontierIdx := GenerateRandomBetween(0, len(frontier))
		chosenFrontier := frontier[randomlyChosenFrontierIdx]

		// Choose one of the Paths 2 cells away to connect to
		availablePathToConnect := m.GetAdjacentPaths(chosenFrontier)
		if len(availablePathToConnect) == 0 {
			panic("Whoah... Somehow we got an invalid frontier")
		}
		randomlyChosenPathToConnect := GenerateRandomBetween(0, len(availablePathToConnect))
		pathToConnect := availablePathToConnect[randomlyChosenPathToConnect]

		// Calculate the cell in-between
		inBetweenX := (chosenFrontier.x + pathToConnect.x) / 2
		inBetweenY := (chosenFrontier.y + pathToConnect.y) / 2

		// Mark the new paths
		m.grid[inBetweenY][inBetweenX] = Path
		m.grid[chosenFrontier.y][chosenFrontier.x] = Path

		// Remove the old chosen frontier from the list
		frontier[randomlyChosenFrontierIdx] = frontier[len(frontier)-1]
		frontier = frontier[:len(frontier)-1]

		// Append new frontiers
		frontier = append(frontier, m.CalculateValidFrontiers(chosenFrontier)...)

		// Paint frontiers
		for _, f := range frontier {
			m.grid[f.y][f.x] = Indicator
		}

		// Sleep for 200ms
		// time.Sleep(200 * time.Millisecond)
	}
}

func (m *Maze) GenerateFirstPrimIteration() []Cell {
	randStartingX := GenerateRandomBetween(2, m.width-2)
	randStartingY := GenerateRandomBetween(2, m.height-2)

	startingCell := Cell{x: randStartingX, y: randStartingY}

	// Just to make sure
	if !m.IsInMaze(startingCell) {
		panic("False random starting location generated!")
	}

	m.grid[randStartingY][randStartingX] = Path

	// Append available starting frontiers
	return m.CalculateValidFrontiers(startingCell)
}

func (m *Maze) CalculateValidFrontiers(cell Cell) []Cell {
	var frontier []Cell

	leftFrontier := Cell{x: cell.x - 2, y: cell.y}
	if m.IsValidFrontier(leftFrontier) {
		frontier = append(frontier, leftFrontier)
	}

	rightFrontier := Cell{x: cell.x + 2, y: cell.y}
	if m.IsValidFrontier(rightFrontier) {
		frontier = append(frontier, rightFrontier)
	}

	topFrontier := Cell{x: cell.x, y: cell.y + 2}
	if m.IsValidFrontier(topFrontier) {
		frontier = append(frontier, topFrontier)
	}

	bottomFrontier := Cell{x: cell.x, y: cell.y - 2}
	if m.IsValidFrontier(bottomFrontier) {
		frontier = append(frontier, bottomFrontier)
	}

	return frontier
}

func (m *Maze) IsValidFrontier(cell Cell) bool {
	isInMaze := m.IsInMaze(cell)
	isAWall := m.grid[cell.y][cell.x] == Wall
	return isInMaze && isAWall
}

func (m *Maze) IsInMaze(cell Cell) bool {
	isValidX := cell.x < m.width-2 && cell.x > 1
	isValidY := cell.y < m.height-2 && cell.y > 1
	return isValidX && isValidY
}

func (m *Maze) GetAdjacentPaths(cell Cell) []Cell {
	var adjacentCells []Cell

	// 2 Cells to the left is a Path
	if cell.x-2 >= 0 {
		if m.grid[cell.y][cell.x-2] == Path {
			adjacentCells = append(adjacentCells, Cell{x: cell.x - 2, y: cell.y})
		}
	}

	// 2 Cells to the right is a Path
	if cell.x+2 < m.width {
		if m.grid[cell.y][cell.x+2] == Path {
			adjacentCells = append(adjacentCells, Cell{x: cell.x + 2, y: cell.y})
		}
	}

	// 2 Cells above is a Path
	if cell.y-2 >= 0 {
		if m.grid[cell.y-2][cell.x] == Path {
			adjacentCells = append(adjacentCells, Cell{x: cell.x, y: cell.y - 2})
		}
	}

	// 2 Cells below is a Path
	if cell.y+2 < m.height {
		if m.grid[cell.y+2][cell.x] == Path {
			adjacentCells = append(adjacentCells, Cell{x: cell.x, y: cell.y + 2})
		}
	}

	return adjacentCells
}

func GenerateRandomBetween(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}
