// Package ui provides GUI rendering functionality for Conway's Game of Life
// using the Fyne toolkit.
package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// InitGUI initializes the graphical user interface grid.
// It creates a grid of rectangle cells and sets up the window content.
// Returns a slice of canvas objects representing each cell in the grid.
func InitGUI(w fyne.Window, gridSize int, cellSize float32) []fyne.CanvasObject {
	color1 := color.RGBA{R: 0, G: 0, B: 0, A: 255} // Black

	var cells []fyne.CanvasObject

	for row := 0; row < gridSize; row++ {
		for col := 0; col < gridSize; col++ {
			rect := canvas.NewRectangle(color1)
			rect.Resize(fyne.NewSize(cellSize, cellSize))
			rect.Move(fyne.NewPos(float32(col)*cellSize, float32(row)*cellSize))
			cells = append(cells, rect)
		}
	}

	content := container.NewWithoutLayout(cells...)
	w.SetContent(content)
	w.Resize(fyne.NewSize(640, 640))
	return cells
}

// UpdateGUI updates the visual state of the grid cells based on the game state.
// Each cell is colored according to whether it is alive (grid[i] == 1) or dead.
func UpdateGUI(cells []fyne.CanvasObject, grid []int, aliveColor, deadColor color.Color) {
	for i := 0; i < len(grid) && i < len(cells); i++ {
		rect := cells[i].(*canvas.Rectangle)
		if grid[i] == 1 {
			rect.FillColor = aliveColor
		} else {
			rect.FillColor = deadColor
		}
		rect.Refresh()
	}
}
