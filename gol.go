package main

import (
	"errors"
	"log"
	"math"
	"math/rand"

	m "github.com/go-gl/mathgl/mgl32"
	p "github.com/jsbento/PhxEngine/renderer/primitives"
	p2d "github.com/jsbento/PhxEngine/renderer/primitives/2d"
)

type GameOfLife struct {
	WindowWidth             int
	WindowHeight            int
	Width                   int
	Height                  int
	Cells                   []bool
	Next                    []bool
	MultiDraw               bool
	CurrentMultiDrawIndexes map[int]bool
}

func NewGameOfLife(width, height, windowWidth, windowHeight int) (*GameOfLife, error) {
	if width <= 0 || height <= 0 || windowWidth <= 0 || windowHeight <= 0 {
		return nil, errors.New(
			"width, height, windowWidth, and windowHeight must be greater than 0",
		)
	}
	if width != height {
		return nil, errors.New("width and height must be equal")
	}
	if width%2 != 0 {
		return nil, errors.New("width must be even")
	}
	if height%2 != 0 {
		return nil, errors.New("height must be even")
	}

	return &GameOfLife{
		WindowWidth:  windowWidth,
		WindowHeight: windowHeight,
		Width:        width,
		Height:       height,
		Cells:        make([]bool, width*height),
		Next:         make([]bool, width*height),
	}, nil
}

func (gol *GameOfLife) ToggleMultiDraw(multiDraw bool) {
	gol.MultiDraw = multiDraw
	if multiDraw {
		gol.CurrentMultiDrawIndexes = make(map[int]bool)
	} else {
		gol.CurrentMultiDrawIndexes = nil
	}
}

func (gol *GameOfLife) Randomize() {
	for i := range gol.Cells {
		gol.Cells[i] = rand.Float64() < 0.25
	}
}

func (gol *GameOfLife) Clear() {
	for i := range gol.Cells {
		gol.Cells[i] = false
	}
}

func (gol *GameOfLife) Update() {
	for y := 0; y < gol.Height; y++ {
		for x := 0; x < gol.Width; x++ {
			liveNeighbors := gol.CountLiveNeighbors(x, y)
			index := y*gol.Width + x
			if gol.Cells[index] {
				gol.Next[index] = liveNeighbors == 2 || liveNeighbors == 3
			} else {
				gol.Next[index] = liveNeighbors == 3
			}
		}
	}
	gol.Cells, gol.Next = gol.Next, gol.Cells
}

func (gol *GameOfLife) ToggleCell(x, y float64) {
	cellX := math.Floor((x + 0.5) * float64(gol.Width))
	cellY := math.Floor((y + 0.5) * float64(gol.Height))

	index := int(cellY)*gol.Width + int(cellX)
	if index < 0 || index >= len(gol.Cells) {
		log.Printf("index out of bounds: %d", index)
		return
	}

	if gol.MultiDraw {
		if toggled := gol.CurrentMultiDrawIndexes[index]; !toggled {
			gol.CurrentMultiDrawIndexes[index] = true
			gol.Cells[index] = !gol.Cells[index]
		}
	} else {
		gol.Cells[index] = !gol.Cells[index]
	}
}

func (gol *GameOfLife) CountLiveNeighbors(x, y int) int {
	count := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := (x+dx+gol.Width)%gol.Width, (y+dy+gol.Height)%gol.Height
			if gol.Cells[ny*gol.Width+nx] {
				count++
			}
		}
	}
	return count
}

func (gol *GameOfLife) Renderables() []p.Renderable2D {
	cellSize := m.Vec2{
		float32(gol.WindowWidth) / float32(gol.Width),
		float32(gol.WindowHeight) / float32(gol.Height),
	}
	offsets := m.Vec2{
		float32(gol.WindowWidth) / 2.0,
		float32(gol.WindowHeight) / 2.0,
	}

	out := []p.Renderable2D{}
	for y := 0; y < gol.Height; y++ {
		for x := 0; x < gol.Width; x++ {
			index := y*gol.Width + x

			if gol.Cells[index] {
				xCoord := (((float32(x) * cellSize.X()) + cellSize.X()/2.0) - offsets.X()) / offsets.X()
				yCoord := (((float32(y) * cellSize.Y()) + cellSize.Y()/2.0) - offsets.Y()) / offsets.Y()

				out = append(out, p.Renderable2D(
					p2d.NewQuad(
						m.Vec2{
							xCoord,
							yCoord,
						},
						m.Vec2{
							(cellSize.X() / float32(gol.WindowWidth)) * 1.1,
							(cellSize.Y() / float32(gol.WindowHeight)) * 1.1,
						},
						0.0,
					)),
				)
			}
		}
	}

	return out
}
