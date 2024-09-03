//TODO: add and delete cells with the mouse
//TODO: zoom in and out with scroll wheel

package main

import (
	"image/color"
	"log"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	grid   [][]bool
	count  int
	paused bool
}

func (g *Game) CountLiveNeighbors(x, y int) int {
	// Initialize a counter for live neighbors
	count := 0

	// Iterate over the 3x3 grid centered on the current cell
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			// Skip the current cell itself
			if i == 0 && j == 0 {
				continue
			}

			// Calculate the coordinates of the neighbor
			ni, nj := x+i, y+j

			// Check if the neighbor is within the grid boundaries
			if ni >= 0 && ni < len(g.grid) && nj >= 0 && nj < len(g.grid[0]) {
				// If the neighbor cell is alive (true), increment the count
				if g.grid[ni][nj] {
					count++
				}
			}
		}
	}

	// Return the total count of live neighbors
	return count
}

// TODO: Function for using mouse to draw / remove cells
func (g *Game) AlterCells() {
}

// TODO: Function to zoom in and out on the grid
func (g *Game) Zoom() {}

func (g *Game) CheckPause() {
	// Pause the game
	//NOTE: this needs some work
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.paused = !g.paused
	}
}

func (g *Game) Update() error {
	g.CheckPause()

	if g.paused != true {
		// Runs runs every 2 frames
		g.count++
		if g.count >= 2 {
			g.count = 0
			//init new grid with same dimensions
			width, height := len(g.grid), len(g.grid[0])
			newGrid := make([][]bool, width)
			for i := range newGrid {
				newGrid[i] = make([]bool, height)
			}

			//Iterate over original grid and check if cell should live or die in next step
			for i := range g.grid {
				for j := range g.grid[i] {
					liveNeighbors := g.CountLiveNeighbors(i, j)

					if g.grid[i][j] {
						//In the case where OG cell is alive
						newGrid[i][j] = liveNeighbors == 2 || liveNeighbors == 3
					} else if !g.grid[i][j] {
						// in the case where the OG cell is dead
						newGrid[i][j] = liveNeighbors == 3
					}
				}
			}
			// set grid to equal our new grid state
			g.grid = newGrid

		}

	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i := range g.grid {
		for j := range g.grid[i] {
			if g.grid[i][j] == true {
				vector.DrawFilledRect(screen, float32(i), float32(j), float32(i*2), float32(j*2), color.White, false)
			} else if g.grid[i][j] == false {
				vector.DrawFilledRect(screen, float32(i), float32(j), 2, 2, color.Black, false)
			}
		}
	}
}

func NewGame(width, height int) *Game {
	// Create a new 2D slice of booleans
	// golang by default initializes bool slices to false
	g := make([][]bool, width)
	for i := range g {
		g[i] = make([]bool, height)
	}

	// init grid elements to randomly be true or false
	for i := range g {
		for j := range g[i] {

			r := rand.IntN(2)
			if r == 0 {
				g[i][j] = false
			} else if r == 1 {
				g[i][j] = true
			}
		}
	}

	// return a pointer to the game with the new grid
	return &Game{
		grid: g,
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 160, 120
}

func main() {
	ebiten.SetWindowSize(960, 720)
	game := NewGame(160, 120)
	ebiten.SetWindowTitle("Conway's Game of Life")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
