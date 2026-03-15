package shaders

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type ShaderProgram struct {
	handle  uint32
	shaders []*Shader
}

func NewShaderProgram(shaderConfigs []ShaderConfig) (*ShaderProgram, error) {
	shaders := []*Shader{}
	for _, cfg := range shaderConfigs {
		shader, err := NewShader(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to compile shader: %v", err)
		}
		shaders = append(shaders, shader)
	}

	program := gl.CreateProgram()
	for _, s := range shaders {
		gl.AttachShader(program, s.handle)
	}
	gl.LinkProgram(program)

	var isLinked int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &isLinked)
	if isLinked == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		logMessage := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(logMessage))

		gl.DeleteProgram(program)
		for _, s := range shaders {
			gl.DeleteShader(s.handle)
		}

		return nil, errors.New(logMessage)
	}

	for _, s := range shaders {
		gl.DetachShader(program, s.handle)
		gl.DeleteShader(s.handle)
	}

	return &ShaderProgram{
		handle:  program,
		shaders: shaders,
	}, nil
}

func (p ShaderProgram) GetHandle() uint32 {
	return p.handle
}

func (p *ShaderProgram) Bind() {
	gl.UseProgram(p.handle)
}

func (p *ShaderProgram) Unbind() {
	gl.UseProgram(0)
}

func (p *ShaderProgram) Cleanup() {
	gl.DeleteProgram(p.handle)
	p.handle = 0
}
