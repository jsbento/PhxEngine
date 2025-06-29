package main

import (
	"flag"
	"log"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	// m "github.com/go-gl/mathgl/mgl32"
	i "github.com/jsbento/PhxEngine/core/init"
	// p "github.com/jsbento/PhxEngine/renderer/primitives"
	// p2d "github.com/jsbento/PhxEngine/renderer/primitives/2d"
	// p3d "github.com/jsbento/PhxEngine/renderer/primitives/3d"

	s "github.com/jsbento/PhxEngine/renderer/shaders"
)

const (
	threshold = 0.15
	fps       = 10

	// width  = 1280
	// height = 720
	width  = 600
	height = 600

	vertexShaderSource = `
    #version 410
    in vec3 vp;
    void main() {
        gl_Position = vec4(vp, 1.0);
    }
	` + "\x00"

	fragmentShaderSource = `
    #version 410
    out vec4 frag_color;
    void main() {
        frag_color = vec4(1.0, 0.0, 0.0, 1.0);
    }
	` + "\x00"
)

func main() {
	isDebug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	runtime.LockOSThread()

	window, err := i.InitGlfw(&i.WindowProps{
		Width:  width,
		Height: height,
		Title:  "Phoenix Engine",
		Debug:  *isDebug,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer glfw.Terminate()

	err = i.InitOpenGL(*isDebug)
	if err != nil {
		log.Fatal("Failed to initialize OpenGL:", err)
	}

	gol, err := NewGameOfLife(10, 10, width, height)
	if err != nil {
		log.Fatal("Failed to create game of life:", err)
	}
	gol.Randomize()

	mousePosX := 0.0
	mousePosY := 0.0

	mousePosToWindowRenderCoords := func(x, y float64) (float64, float64) {
		renderX := (x - float64(width)/2.0) / float64(width)
		renderY := -1.0 * (y - float64(height)/2.0) / float64(height)
		return renderX, renderY
	}

	window.SetMouseButtonCallback(
		func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
			switch action {
			case glfw.Press:
				gol.ToggleMultiDraw(true)
				x, y := mousePosToWindowRenderCoords(mousePosX, mousePosY)
				if button == glfw.MouseButton1 {
					gol.ToggleCell(x, y)
				}
			case glfw.Release:
				gol.ToggleMultiDraw(false)
			}
		},
	)

	paused := false
	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if key == glfw.KeySpace && action == glfw.Press {
			paused = !paused
		}
		if key == glfw.KeyEscape && action == glfw.Press {
			gol.Clear()
		}
	})

	window.SetCursorPosCallback(func(w *glfw.Window, xpos, ypos float64) {
		mousePosX = xpos
		mousePosY = ypos

		if gol.MultiDraw {
			x, y := mousePosToWindowRenderCoords(xpos, ypos)
			gol.ToggleCell(x, y)
		}
	})

	program := createProgram()
	for !window.ShouldClose() {
		t := time.Now()
		draw(window, program, gol)
		if !paused {
			gol.Update()
		}
		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
	}
}

func createProgram() uint32 {
	vertexShader, err := s.CompileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		log.Fatalf("Error compiling shader: %v", err)
	}
	defer gl.DeleteShader(vertexShader)

	fragmentShader, err := s.CompileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		log.Fatalf("Error compiling shader: %v", err)
	}
	defer gl.DeleteShader(fragmentShader)

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	return program
}

func draw(window *glfw.Window, program uint32, gol *GameOfLife) {
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	for _, r := range gol.Renderables() {
		r.Draw()
	}

	// renderables := []p.Renderable2D{}
	// renderables3D := []p.Renderable3D{}

	// renderables = append(renderables, p.Renderable2D(p2d.NewQuad(
	// 	m.Vec2{0.5, 0.5},
	// 	m.Vec2{1.0, 1.0},
	// 	m.DegToRad(45.0),
	// )))

	// renderables = append(renderables, p.Renderable2D(p2d.NewCircle(
	// 	m.Vec2{0.7, -0.7},
	// 	m.Vec2{1.0, 1.0},
	// 	0.3,
	// 	50,
	// )))

	// renderables3D = append(renderables3D, p.Renderable3D(p3d.NewCuboid(
	// 	m.Vec3{0.0, -0.5, 0.0},
	// 	m.Vec3{0.5, 0.5, 0.5},
	// 	m.Vec3{0.0, 0.0, 0.0},
	// )))

	// for _, renderable := range renderables {
	// 	renderable.Draw()
	// }
	// for _, renderable := range renderables3D {
	// 	renderable.Draw()
	// }

	// line := p2d.NewLine(
	// 	m.Vec2{-0.75, -0.5},
	// 	m.Vec2{-0.25, 0.75},
	// 	2.0,
	// )
	// line.Draw()

	glfw.PollEvents()
	window.SwapBuffers()
}
