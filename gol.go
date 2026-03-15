package main

import (
	"errors"
	"log"
	"math"
	"math/rand"

	"github.com/go-gl/gl/v4.6-core/gl"
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
	CurrentMultiDrawIndexes []bool
	cellSize                m.Vec2
	offsets                 m.Vec2
	cellScale               m.Vec2
	renderInitialized       bool
	vao                     uint32
	vbo                     uint32
	instanceVbo             uint32
	instancePositions       []float32
	instanceCount           int32
	instanceCapacityBytes   int
}

func NewGameOfLife(width, height int) (*GameOfLife, error) {
	if width <= 0 || height <= 0 {
		return nil, errors.New(
			"width and height must be greater than 0",
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

	gol := &GameOfLife{
		Width:  width,
		Height: height,
		Cells:  make([]bool, width*height),
		Next:   make([]bool, width*height),
	}
	return gol, nil
}

func (gol *GameOfLife) ToggleMultiDraw(multiDraw bool) {
	gol.MultiDraw = multiDraw
	if multiDraw {
		gol.CurrentMultiDrawIndexes = make([]bool, len(gol.Cells))
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
	if cellX < 0 {
		cellX = 0
	} else if cellX >= float64(gol.Width) {
		cellX = float64(gol.Width - 1)
	}
	if cellY < 0 {
		cellY = 0
	} else if cellY >= float64(gol.Height) {
		cellY = float64(gol.Height - 1)
	}

	index := int(cellY)*gol.Width + int(cellX)
	if index < 0 || index >= len(gol.Cells) {
		log.Printf("index out of bounds: %d", index)
		return
	}

	if gol.MultiDraw {
		if gol.CurrentMultiDrawIndexes == nil ||
			len(gol.CurrentMultiDrawIndexes) != len(gol.Cells) {
			gol.CurrentMultiDrawIndexes = make([]bool, len(gol.Cells))
		}
		if !gol.CurrentMultiDrawIndexes[index] {
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
	cellSize := gol.cellSize
	offsets := gol.offsets

	out := make([]p.Renderable2D, 0, gol.Width*gol.Height)
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
							gol.cellScale.X(),
							gol.cellScale.Y(),
						},
						0.0,
					)),
				)
			}
		}
	}

	return out
}

func (gol *GameOfLife) InitRender() {
	if gol.renderInitialized {
		return
	}
	gol.recalculateRenderMetrics()
	baseQuad := []float32{
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		0.5, 0.5, 0.0,
		0.5, 0.5, 0.0,
		-0.5, 0.5, 0.0,
		-0.5, -0.5, 0.0,
	}

	gl.GenVertexArrays(1, &gol.vao)
	gl.BindVertexArray(gol.vao)

	gl.GenBuffers(1, &gol.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, gol.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(baseQuad)*4, gl.Ptr(baseQuad), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)

	gl.GenBuffers(1, &gol.instanceVbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, gol.instanceVbo)
	gol.instanceCapacityBytes = gol.Width * gol.Height * 2 * 4
	gl.BufferData(gl.ARRAY_BUFFER, gol.instanceCapacityBytes, nil, gl.DYNAMIC_DRAW)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribDivisor(1, 1)

	gl.BindVertexArray(0)
	gol.renderInitialized = true
}

func (gol *GameOfLife) DrawInstances(program uint32) {
	gol.InitRender()
	gol.updateInstanceData()

	scaleLoc := gl.GetUniformLocation(program, gl.Str("uScale\x00"))
	if scaleLoc >= 0 {
		gl.Uniform2f(scaleLoc, gol.cellScale.X(), gol.cellScale.Y())
	}

	gl.BindVertexArray(gol.vao)
	gl.DrawArraysInstanced(gl.TRIANGLES, 0, 6, gol.instanceCount)
	gl.BindVertexArray(0)
}

func (gol *GameOfLife) updateInstanceData() {
	capacity := gol.Width * gol.Height * 2
	if cap(gol.instancePositions) < capacity {
		gol.instancePositions = make([]float32, 0, capacity)
	} else {
		gol.instancePositions = gol.instancePositions[:0]
	}

	cellSize := gol.cellSize
	offsets := gol.offsets
	for y := 0; y < gol.Height; y++ {
		for x := 0; x < gol.Width; x++ {
			index := y*gol.Width + x
			if !gol.Cells[index] {
				continue
			}
			xCoord := (((float32(x) * cellSize.X()) + cellSize.X()/2.0) - offsets.X()) / offsets.X()
			yCoord := (((float32(y) * cellSize.Y()) + cellSize.Y()/2.0) - offsets.Y()) / offsets.Y()
			gol.instancePositions = append(gol.instancePositions, xCoord, yCoord)
		}
	}

	gol.instanceCount = int32(len(gol.instancePositions) / 2)
	if len(gol.instancePositions) == 0 {
		return
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, gol.instanceVbo)
	requiredBytes := len(gol.instancePositions) * 4
	if requiredBytes > gol.instanceCapacityBytes {
		gol.instanceCapacityBytes = requiredBytes
		gl.BufferData(gl.ARRAY_BUFFER, gol.instanceCapacityBytes, nil, gl.DYNAMIC_DRAW)
	}
	gl.BufferSubData(
		gl.ARRAY_BUFFER,
		0,
		len(gol.instancePositions)*4,
		gl.Ptr(gol.instancePositions),
	)
}

func (gol *GameOfLife) recalculateRenderMetrics() {
	gol.cellSize = m.Vec2{
		float32(gol.WindowWidth) / float32(gol.Width),
		float32(gol.WindowHeight) / float32(gol.Height),
	}
	gol.offsets = m.Vec2{
		float32(gol.WindowWidth) / 2.0,
		float32(gol.WindowHeight) / 2.0,
	}
	gol.cellScale = m.Vec2{
		(gol.cellSize.X() / float32(gol.WindowWidth)) * 1.1,
		(gol.cellSize.Y() / float32(gol.WindowHeight)) * 1.1,
	}
}

func (gol *GameOfLife) Destroy() {
	if !gol.renderInitialized {
		return
	}
	gl.DeleteBuffers(1, &gol.vbo)
	gl.DeleteBuffers(1, &gol.instanceVbo)
	gl.DeleteVertexArrays(1, &gol.vao)
	gol.vbo = 0
	gol.instanceVbo = 0
	gol.vao = 0
	gol.renderInitialized = false
}

func (gol *GameOfLife) OnResize(windowWidth, windowHeight int) {
	gol.SetWindowSize(windowWidth, windowHeight)
}

func (gol *GameOfLife) SetWindowSize(windowWidth, windowHeight int) {
	if windowWidth <= 0 || windowHeight <= 0 {
		return
	}
	gol.WindowWidth = windowWidth
	gol.WindowHeight = windowHeight
	gol.recalculateRenderMetrics()
	if gol.renderInitialized {
		requiredBytes := gol.Width * gol.Height * 2 * 4
		if requiredBytes > gol.instanceCapacityBytes {
			gol.instanceCapacityBytes = requiredBytes
			gl.BindBuffer(gl.ARRAY_BUFFER, gol.instanceVbo)
			gl.BufferData(gl.ARRAY_BUFFER, gol.instanceCapacityBytes, nil, gl.DYNAMIC_DRAW)
		}
	}
}
