package primitives

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	m "github.com/go-gl/mathgl/mgl32"
	"github.com/jsbento/PhxEngine/math"
)

type Circle struct {
	renderId    uint32
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
	c.UpdateTransform()
	return
}

func (c *Circle) Draw() {
	circle := m.Circle(
		c.Radius*c.Scalar.X(),
		c.Radius*c.Scalar.Y(),
		c.NumSlices,
	)
	vertices := []float32{}
	for _, vertex := range circle {
		v := c.transform.Mul4x1(vertex.Vec4(0.0, 1.0))
		vertices = append(vertices, v.X(), v.Y(), v.Z())
	}

	gl.GenBuffers(1, &c.renderId)
	gl.BindBuffer(gl.ARRAY_BUFFER, c.renderId)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, c.renderId)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertices)))
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
}
