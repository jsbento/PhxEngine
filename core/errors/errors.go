package errors

import (
	"errors"
	"log"

	"github.com/go-gl/gl/v4.6-core/gl"

	"github.com/jsbento/PhxEngine/utils"
)

func CheckGLError() error {
	if err := gl.GetError(); err != gl.NO_ERROR {
		errStr := gl.GoStr(gl.GetString(err))

		log.Printf("OpenGL Error: %s", utils.GetCallerFuncName())
		log.Printf("OpenGL Error: %v", err)
		log.Printf("OpenGL Error: %s", errStr)
		log.Printf("OpenGL Error: %d", err)

		return errors.New(errStr)
	}

	return nil
}
