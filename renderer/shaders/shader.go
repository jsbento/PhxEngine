package shaders

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
)

func CompileShader(source string, shaderType uint32) (uint32, error) {
	log.Printf("OpenGL shader compile start: type=%s\n", shaderTypeString(shaderType))
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		logMessage := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(logMessage))

		log.Printf(
			"OpenGL shader compile failed: type=%s error=%s\n",
			shaderTypeString(shaderType),
			logMessage,
		)
		return 0, fmt.Errorf("failed to compile %v: %v", source, logMessage)
	}

	log.Printf("OpenGL shader compile success: type=%s\n", shaderTypeString(shaderType))
	return shader, nil
}

func shaderTypeString(shaderType uint32) string {
	switch shaderType {
	case gl.VERTEX_SHADER:
		return "vertex"
	case gl.FRAGMENT_SHADER:
		return "fragment"
	case gl.GEOMETRY_SHADER:
		return "geometry"
	case gl.TESS_CONTROL_SHADER:
		return "tess_control"
	case gl.TESS_EVALUATION_SHADER:
		return "tess_evaluation"
	case gl.COMPUTE_SHADER:
		return "compute"
	default:
		return "unknown"
	}
}
