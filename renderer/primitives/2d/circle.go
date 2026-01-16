package primitives

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	m "github.com/go-gl/mathgl/mgl32"
	"github.com/jsbento/PhxEngine/math"
)

type Circle struct {
	renderId    uint32
	vao         uint32
	vertexCount int32
	baseCircle  []m.Vec2
	vertexData  []float32
	buffersInit bool
	dirty       bool
	transform   m.Mat4
	Translation m.Vec2
	Scalar      m.Vec2
	Rotation    float32
	Radius      float32
	NumSlices   int
}

func NewCircle(
	translation m.Vec2,
	scale m.Vec2,
	radius float32,
	numSlices int) (c *Circle) {
	c = &Circle{
		Translation: translation,
		Scalar:      scale,
		Rotation:    0.0,
		Radius:      radius,
		NumSlices:   numSlices,
	}
	c.baseCircle = m.Circle(
		c.Radius,
		c.Radius,
		c.NumSlices,
	)
	c.UpdateTransform()
	return
}

func (c *Circle) Draw() {
	c.ensureBuffers()
	c.updateVertexData()
	gl.BindVertexArray(c.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, c.vertexCount)
	gl.BindVertexArray(0)
}

func (c *Circle) Translate(pos m.Vec2) {
	c.Translation = pos
	c.UpdateTransform()
}

func (c *Circle) Scale(scale m.Vec2) {
	c.Scalar = scale
	c.UpdateTransform()
}

func (c *Circle) Rotate(rotation float32) {
	c.Rotation = rotation
	c.UpdateTransform()
}

func (c *Circle) UpdateTransform() {
	c.transform = math.MakeTransformationMatrix(
		c.Scalar.Vec3(1.0),
		c.Rotation,
		m.Vec3{0.0, 0.0, 1.0},
		c.Translation.Vec3(0.0),
	)
	c.dirty = true
}

func (c *Circle) Destroy() {
	if !c.buffersInit {
		return
	}
	gl.DeleteBuffers(1, &c.renderId)
	gl.DeleteVertexArrays(1, &c.vao)
	c.renderId = 0
	c.vao = 0
	c.buffersInit = false
}

func (c *Circle) ensureBuffers() {
	if c.buffersInit {
		return
	}
	gl.GenVertexArrays(1, &c.vao)
	gl.GenBuffers(1, &c.renderId)

	gl.BindVertexArray(c.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, c.renderId)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)
	gl.BindVertexArray(0)

	c.buffersInit = true
}

func (c *Circle) updateVertexData() {
	if !c.dirty {
		return
	}
	if c.baseCircle == nil {
		c.baseCircle = m.Circle(
			c.Radius,
			c.Radius,
			c.NumSlices,
		)
	}
	required := len(c.baseCircle) * 3
	if cap(c.vertexData) < required {
		c.vertexData = make([]float32, 0, required)
	} else {
		c.vertexData = c.vertexData[:0]
	}
	for _, vertex := range c.baseCircle {
		v := c.transform.Mul4x1(vertex.Vec4(0.0, 1.0))
		c.vertexData = append(c.vertexData, v.X(), v.Y(), v.Z())
	}
	c.vertexCount = int32(len(c.vertexData) / 3)

	gl.BindBuffer(gl.ARRAY_BUFFER, c.renderId)
	gl.BufferData(gl.ARRAY_BUFFER, len(c.vertexData)*4, gl.Ptr(c.vertexData), gl.DYNAMIC_DRAW)
	c.dirty = false
}
