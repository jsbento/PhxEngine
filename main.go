package main

import (
	"log"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	i "github.com/jsbento/PhxEngine/core/init"
	p "github.com/jsbento/PhxEngine/renderer/primitives"
	p2d "github.com/jsbento/PhxEngine/renderer/primitives/2d"
	p3d "github.com/jsbento/PhxEngine/renderer/primitives/3d"
	s "github.com/jsbento/PhxEngine/renderer/shaders"

	m "github.com/go-gl/mathgl/mgl32"
)

const (
	threshold = 0.15
	fps       = 60

	width  = 1280
	height = 720

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
        frag_color = vec4(1, 0, 1, 1);
    }
	` + "\x00"
)

func main() {
	runtime.LockOSThread()

	window, err := i.InitGlfw(&i.WindowProps{
		Width:  width,
		Height: height,
		Title:  "Phoenix Engine",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer glfw.Terminate()

	program := initOpenGL()
	for !window.ShouldClose() {
		t := time.Now()
		draw(window, program)
		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
	}
}

func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		log.Fatal("Failed to initialize gl:", err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version:", version)

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

func draw(window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	renderables := []p.Renderable2D{}
	renderables3D := []p.Renderable3D{}

	renderables = append(renderables, p.Renderable2D(p2d.NewQuad(
		m.Vec2{0.5, 0.5},
		m.Vec2{1.0, 1.0},
		m.DegToRad(45.0),
	)))

	renderables = append(renderables, p.Renderable2D(p2d.NewCircle(
		m.Vec2{0.7, -0.7},
		m.Vec2{1.0, 1.0},
		0.3,
		50,
	)))

	renderables3D = append(renderables3D, p.Renderable3D(p3d.NewCuboid(
		m.Vec3{0.0, -0.5, 0.0},
		m.Vec3{0.5, 0.5, 0.5},
		m.Vec3{0.0, 0.0, 0.0},
	)))

	for _, renderable := range renderables {
		renderable.Draw()
	}
	for _, renderable := range renderables3D {
		renderable.Draw()
	}

	line := p2d.NewLine(
		m.Vec2{-0.75, -0.5},
		m.Vec2{-0.25, 0.75},
		2.0,
	)
	line.Draw()

	glfw.PollEvents()
	window.SwapBuffers()
}
