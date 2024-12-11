package main

import (
	"strings"
)

const (
	Wall  = '#'
	Path  = ' '
	Start = 'S'
	End   = 'E'
)

type Cell struct {
	x, y int
}

type Maze struct {
	height, width int
	grid          [][]rune
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
	var sb strings.Builder
	for _, row := range m.grid {
		sb.WriteString(string(row) + "\n")
	}
	return sb.String()
}

// func GenerateMaze() {
//
// }
