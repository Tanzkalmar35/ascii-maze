package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.org/x/term"
)

func GenerateExampleMaze(width, height int) *Maze {

	// Initialize the maze
	maze := DisplayMaze(width, height-1)

	maze.grid[1][1] = Start // Set start point
	maze.grid[8][8] = End   // Set end point
	maze.grid[1][2] = Path  // Create a path
	maze.grid[1][3] = Path  // Create a path
	maze.grid[2][3] = Path  // Create a path
	maze.grid[3][3] = Path  // Create a path
	maze.grid[4][3] = Path  // Create a path
	maze.grid[5][3] = Path  // Create a path
	maze.grid[6][3] = Path  // Create a path
	maze.grid[7][3] = Path  // Create a path
	maze.grid[7][4] = Path  // Create a path
	maze.grid[7][5] = Path  // Create a path
	maze.grid[7][6] = Path  // Create a path
	maze.grid[7][7] = Path  // Create a path
	maze.grid[8][7] = Path  // Create a path

	return maze
}

func SetupApplication() {
}

func main() {

	width, height, err := term.GetSize(int(os.Stdout.Fd()))

	if err != nil {
		panic("Could not get terminal size")
	}

	maze := GenerateExampleMaze(width, height)
	app := tview.NewApplication()

	infoView := tview.NewTextView().SetText("Press 'q' to quit.").SetTextAlign(tview.AlignCenter)
	mazeView := tview.NewTextView().SetText(maze.Render()).SetTextAlign(tview.AlignLeft)

	mazeView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Rune() == 'q' {
			app.Stop()
		}
		return event
	})

	flexView := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(infoView, 1, 0, false).
		AddItem(mazeView, 0, 1, true)

	if err := app.SetRoot(flexView, true).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running application: &v\n", err)
		os.Exit(1)
	}
}
