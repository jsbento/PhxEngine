package init

import (
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
)

func InitOpenGL(debug bool) error {
	if err := gl.Init(); err != nil {
		return err
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	println("OpenGL version:", version)

	if debug {
		gl.Enable(gl.DEBUG_OUTPUT)
		gl.DebugMessageCallback(
			func(source, gltype, id uint32, severity uint32, length int32, message string, userParam unsafe.Pointer) {
				println("OpenGL Debug:", message)
				if severity == gl.DEBUG_SEVERITY_HIGH {
					println("OpenGL Error:", message)
					panic("OpenGL Error")
				}
				if severity == gl.DEBUG_SEVERITY_MEDIUM {
					println("OpenGL Warning:", message)
				}
				if severity == gl.DEBUG_SEVERITY_LOW {
					println("OpenGL Info:", message)
				}
				if severity == gl.DEBUG_SEVERITY_NOTIFICATION {
					println("OpenGL Notification:", message)
				}
			},
			nil,
		)
		gl.DebugMessageControl(gl.DONT_CARE, gl.DONT_CARE, gl.DONT_CARE, 0, nil, true)
	}

	return nil
}
