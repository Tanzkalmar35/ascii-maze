package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func GenerateMaze(width, height int, prim bool) *Maze {
	// Initialize the maze
	maze := DisplayMaze(width, height)

	// Generate a random perfect maze using randomized prim's algorighm
	if prim {
		maze.GenerateRandomPrimMaze()
	}

	currentMaze = maze

	return maze
}

func main() {
	// Create the tea gui
	p := tea.NewProgram(model{quit: false}, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
