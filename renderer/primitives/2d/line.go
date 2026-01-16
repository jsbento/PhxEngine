package primitives

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	m "github.com/go-gl/mathgl/mgl32"
)

type Line struct {
	renderId     uint32
	vao          uint32
	vertexCount  int32
	vertexData   []float32
	buffersInit  bool
	dirty        bool
	cachedStart  m.Vec2
	cachedEnd    m.Vec2
	cachedWeight float32
	Start        m.Vec2
	End          m.Vec2
	Thickness    float32
}

func NewLine(start, end m.Vec2, thickness float32) *Line {
	return &Line{
		Start:     start,
		End:       end,
		Thickness: thickness,
	}
}

func (l *Line) Draw() {
	l.ensureBuffers()
	l.updateVertexData()
	gl.BindVertexArray(l.vao)
	gl.LineWidth(l.Thickness)
	gl.DrawArrays(gl.LINES, 0, l.vertexCount)
	gl.LineWidth(1.0)
	gl.BindVertexArray(0)
}

func (l *Line) ensureBuffers() {
	if l.buffersInit {
		return
	}
	gl.GenVertexArrays(1, &l.vao)
	gl.GenBuffers(1, &l.renderId)

	gl.BindVertexArray(l.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, l.renderId)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)
	gl.BindVertexArray(0)

	l.buffersInit = true
	l.dirty = true
}

func (l *Line) updateVertexData() {
	if l.Start != l.cachedStart || l.End != l.cachedEnd || l.Thickness != l.cachedWeight {
		l.dirty = true
	}
	if !l.dirty {
		return
	}
	required := 6
	if cap(l.vertexData) < required {
		l.vertexData = make([]float32, 0, required)
	} else {
		l.vertexData = l.vertexData[:0]
	}
	l.vertexData = append(l.vertexData,
		l.Start.X(), l.Start.Y(), 0.0,
		l.End.X(), l.End.Y(), 0.0,
	)
	l.vertexCount = int32(len(l.vertexData) / 3)
	l.cachedStart = l.Start
	l.cachedEnd = l.End
	l.cachedWeight = l.Thickness

	gl.BindBuffer(gl.ARRAY_BUFFER, l.renderId)
	gl.BufferData(gl.ARRAY_BUFFER, len(l.vertexData)*4, gl.Ptr(l.vertexData), gl.DYNAMIC_DRAW)
	l.dirty = false
}

func (l *Line) Destroy() {
	if !l.buffersInit {
		return
	}
	gl.DeleteBuffers(1, &l.renderId)
	gl.DeleteVertexArrays(1, &l.vao)
	l.renderId = 0
	l.vao = 0
	l.buffersInit = false
}
