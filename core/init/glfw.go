package init

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

type WindowProps struct {
	Width  int
	Height int
	Title  string
}

func InitGlfw(windowProps *WindowProps) (*glfw.Window, error) {
	if err := glfw.Init(); err != nil {
		return nil, err
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(
		windowProps.Width,
		windowProps.Height,
		windowProps.Title,
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}
	window.MakeContextCurrent()

	return window, nil
}
