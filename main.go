//TODO: zoom in and out with scroll wheel

package main

import (
	"image/color"
	"log"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const WIDTH int = 192
const HEIGHT int = 144
const SCALE int = 12
const TICK_SPEED int = 2

type Game struct {
	grid     [][]bool
	count    int
	paused   bool
	editMode string
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

// TODO: Function to zoom in and out on the grid
func (g *Game) Zoom() {}

func (g *Game) HandlePause() {
	// Pause the game
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.paused = !g.paused
	}
}

func (g *Game) Update() error {
	// listen for pause (spacebar)
	g.HandlePause()

	//Main loop
	if g.paused != true {
		// Runs runs every TICK_SPEED ticks
		g.count++
		if g.count >= TICK_SPEED {
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
	} else if g.paused {
		// allow user to swap between pencil & eraser mode when paused
		if g.paused && inpututil.IsKeyJustPressed(ebiten.KeyE) {
			switch g.editMode {
			case "pencil":
				g.editMode = "eraser"
			case "eraser":
				g.editMode = "pencil"
			}
		}

		// Allow user to draw their own cells when paused
		x, y := ebiten.CursorPosition()
		if g.paused && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) == true {
			switch g.editMode {
			case "pencil":
				g.grid[x][y] = true
			case "eraser":
				g.grid[x][y] = false
			}
		}

	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i := range g.grid {
		for j := range g.grid[i] {
			if g.grid[i][j] == true {
				vector.DrawFilledRect(screen, float32(i), float32(j), float32(i*SCALE), float32(j*SCALE), color.RGBA{128, 128, 128, 255}, false)
			} else if g.grid[i][j] == false {
				vector.DrawFilledRect(screen, float32(i), float32(j), float32(i*SCALE), float32(j*SCALE), color.Black, false)
			}
		}
	}

	// TODO: add zoom

	if g.paused {
		switch g.editMode {
		case "pencil":
			ebitenutil.DebugPrint(screen, "Paused \nClick & drag to add new cells \nPress E to switch to eraser")
		case "eraser":
			ebitenutil.DebugPrint(screen, "Paused \nClick & drag to remove cells \nPress E to switch to pencil")
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
		grid:     g,
		editMode: "pencil",
		paused:   false,
		count:    0,
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WIDTH, HEIGHT
}

func main() {
	ebiten.SetWindowSize(WIDTH*(SCALE/2), HEIGHT*(SCALE/2))
	game := NewGame(WIDTH, HEIGHT)
	ebiten.SetWindowTitle("Conway's Game of Life")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
