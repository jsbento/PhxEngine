package errors

import (
	"log"

	"github.com/go-gl/gl/v4.6-core/gl"
)

func CheckGLError(caller string) {
	if err := gl.GetError(); err != gl.NO_ERROR {
		log.Printf("OpenGL Error: %s", caller)
		log.Printf("OpenGL Error: %v", err)
		log.Printf("OpenGL Error: %s", gl.GoStr(gl.GetString(err)))
		log.Printf("OpenGL Error: %d", err)
	}
}
