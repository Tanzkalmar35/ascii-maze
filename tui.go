package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

type model struct {
	quit bool
	maze Maze
}

func (m model) Init() tea.Cmd {
	// Generate the maze on startup
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic("Could not get terminal size")
	}

	currentMaze = GenerateMaze(width, height-1, false) // Initialize the maze
	// go currentMaze.GenerateRandomPrimMaze()  // Start generating the maze in a goroutine

	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return model{quit: true}, tea.Quit
		case "w":
			// Regenerate maze
			ReGenerateMaze(true)
			SolvingStrategy{}.AStarAlgo(*currentMaze)
		}
	}
	return m, nil
}

func (m model) View() string {
	if currentMaze != nil {
		return currentMaze.Render()
	}

	return "Maze still loading..."
}

func ReGenerateMaze(prim bool) string {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))

	if err != nil {
		panic("Could not get terminal size")
	}

	maze := GenerateMaze(width, height-1, prim)
	currentMaze = maze
	return maze.Render()
}
