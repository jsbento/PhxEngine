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
	fps = 10

	vertexShaderSource = `
    #version 410
    layout (location = 0) in vec3 vp;
    layout (location = 1) in vec2 instancePos;
    uniform vec2 uScale;
    void main() {
        vec2 pos = (vp.xy * uScale) + instancePos;
        gl_Position = vec4(pos, vp.z, 1.0);
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

	initialWindowWidth := 600
	initialWindowHeight := 600
	window, err := i.InitGlfw(&i.WindowProps{
		Width:  initialWindowWidth,
		Height: initialWindowHeight,
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
	framebufferWidth, framebufferHeight := window.GetFramebufferSize()
	if framebufferWidth > 0 && framebufferHeight > 0 {
		gl.Viewport(0, 0, int32(framebufferWidth), int32(framebufferHeight))
	}

	gol, err := NewGameOfLife(100, 100)
	if err != nil {
		log.Fatal("Failed to create game of life:", err)
	}
	windowWidth, windowHeight := window.GetSize()
	gol.SetWindowSize(windowWidth, windowHeight)
	gol.Randomize()
	gol.InitRender()
	defer gol.Destroy()
	mousePosX := 0.0
	mousePosY := 0.0

	mousePosToWindowRenderCoords := func(x, y float64) (float64, float64) {
		renderX := (x - float64(windowWidth)/2.0) / float64(windowWidth)
		renderY := -1.0 * (y - float64(windowHeight)/2.0) / float64(windowHeight)
		return renderX, renderY
	}

	window.SetFramebufferSizeCallback(func(w *glfw.Window, fbWidth, fbHeight int) {
		if fbWidth <= 0 || fbHeight <= 0 {
			return
		}
		gl.Viewport(0, 0, int32(fbWidth), int32(fbHeight))
		windowWidth, windowHeight = w.GetSize()
		gol.OnResize(windowWidth, windowHeight)
	})

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
	window.SetKeyCallback(
		func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
			if key == glfw.KeySpace && action == glfw.Press {
				paused = !paused
			}
			if key == glfw.KeyEscape && action == glfw.Press {
				gol.Clear()
			}
		},
	)

	window.SetCursorPosCallback(func(w *glfw.Window, xpos, ypos float64) {
		mousePosX = xpos
		mousePosY = ypos

		if gol.MultiDraw {
			x, y := mousePosToWindowRenderCoords(xpos, ypos)
			gol.ToggleCell(x, y)
		}
	})

	program, err := s.NewShaderProgram([]s.ShaderConfig{
		{
			Source:     vertexShaderSource,
			ShaderType: gl.VERTEX_SHADER,
		},
		{
			Source:     fragmentShaderSource,
			ShaderType: gl.FRAGMENT_SHADER,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer program.Cleanup()
	for !window.ShouldClose() {
		t := time.Now()
		draw(window, program, gol)
		if !paused {
			gol.Update()
		}
		frameBudget := time.Second / time.Duration(fps)
		elapsed := time.Since(t)
		if remaining := frameBudget - elapsed; remaining > 0 {
			time.Sleep(remaining)
		}
	}
}

func draw(window *glfw.Window, program *s.ShaderProgram, gol *GameOfLife) {
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program.GetHandle())
	gol.DrawInstances(program.GetHandle())

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
	// 	renderable.Destroy()
	// }
	// for _, renderable := range renderables3D {
	// 	renderable.Draw()
	// 	renderable.Destroy()
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
