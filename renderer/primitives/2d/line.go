package primitives

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	m "github.com/go-gl/mathgl/mgl32"
)

type Line struct {
	renderId  uint32
	Start     m.Vec2
	End       m.Vec2
	Thickness float32
}

func NewLine(start, end m.Vec2, thickness float32) *Line {
	return &Line{
		Start:     start,
		End:       end,
		Thickness: thickness,
	}
}

func (l *Line) Draw() {
	vertices := []float32{
		l.Start.X(), l.Start.Y(), 0.0,
		l.End.X(), l.End.Y(), 0.0,
	}

	gl.GenBuffers(1, &l.renderId)
	gl.BindBuffer(gl.ARRAY_BUFFER, l.renderId)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, l.renderId)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	gl.LineWidth(l.Thickness)
	gl.DrawArrays(gl.LINES, 0, int32(len(vertices)))
	gl.LineWidth(1.0)
}
