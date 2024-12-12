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
		}
	}
	return m, nil
}

func (m model) View() string {
	if currentMaze != nil {
		return currentMaze.Render()
	}

	return "Press 'q' to quit!\n" + ReGenerateMaze(false)
}

func ReGenerateMaze(prim bool) string {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))

	if err != nil {
		panic("Could not get terminal size")
	}

	maze := GenerateMaze(width, height, prim)
	return maze.Render()
}
