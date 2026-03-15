package shaders

import (
	"errors"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type ShaderConfig struct {
	Source     string
	ShaderType uint32
}

type Shader struct {
	handle        uint32
	shaderType    uint32
	shaderTypeStr string
	source        string
}

func NewShader(config ShaderConfig) (*Shader, error) {
	shader := &Shader{
		source:        config.Source,
		shaderType:    config.ShaderType,
		shaderTypeStr: shaderTypeString(config.ShaderType),
	}
	if err := shader.compile(); err != nil {
		return nil, err
	}
	return shader, nil
}

func (s Shader) GetHandle() uint32 {
	return s.handle
}

func (s *Shader) compile() error {
	shader := gl.CreateShader(s.shaderType)

	source, free := gl.Strs(s.source)
	gl.ShaderSource(shader, 1, source, nil)
	free()
	gl.CompileShader(shader)

	var compileStatus int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &compileStatus)
	if compileStatus == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		logMessage := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(logMessage))
		gl.DeleteShader(shader)
		return errors.New(logMessage)
	}
	s.handle = shader
	return nil
}

func (s *Shader) Cleanup() {
	gl.DeleteShader(s.handle)
	s.handle = 0
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
