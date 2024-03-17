package primitives

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	m "github.com/go-gl/mathgl/mgl32"
	"github.com/jsbento/PhxEngine/math"
	"github.com/jsbento/PhxEngine/renderer/utils"
)

type Cuboid struct {
	renderId    uint32
	transform   m.Mat4
	Translation m.Vec3
	Scalar      m.Vec3
	Rotation    m.Vec3
}

func NewCuboid(
	translation m.Vec3,
	scale m.Vec3,
	rotation m.Vec3,
) (c *Cuboid) {
	c = &Cuboid{
		Translation: translation,
		Scalar:      scale,
		Rotation:    rotation,
	}
	c.UpdateTransform()
	return
}

func (c *Cuboid) Draw() {
	cuboid := []m.Vec4{
		// front face
		{-0.5, -0.5, 0.5, 1.0},
		{0.5, -0.5, 0.5, 1.0},
		{0.5, 0.5, 0.5, 1.0},
		{0.5, 0.5, 0.5, 1.0},
		{-0.5, 0.5, 0.5, 1.0},
		{-0.5, -0.5, 0.5, 1.0},
		// back face
		{-0.5, -0.5, -0.5, 1.0},
		{0.5, -0.5, -0.5, 1.0},
		{0.5, 0.5, -0.5, 1.0},
		{0.5, 0.5, -0.5, 1.0},
		{-0.5, 0.5, -0.5, 1.0},
		{-0.5, -0.5, -0.5, 1.0},
		// left face
		{-0.5, -0.5, -0.5, 1.0},
		{-0.5, -0.5, 0.5, 1.0},
		{-0.5, 0.5, 0.5, 1.0},
		{-0.5, 0.5, 0.5, 1.0},
		{-0.5, 0.5, -0.5, 1.0},
		{-0.5, -0.5, -0.5, 1.0},
		// right face
		{0.5, -0.5, -0.5, 1.0},
		{0.5, -0.5, 0.5, 1.0},
		{0.5, 0.5, 0.5, 1.0},
		{0.5, 0.5, 0.5, 1.0},
		{0.5, 0.5, -0.5, 1.0},
		{0.5, -0.5, -0.5, 1.0},
		// top face
		{-0.5, 0.5, -0.5, 1.0},
		{0.5, 0.5, -0.5, 1.0},
		{0.5, 0.5, 0.5, 1.0},
		{0.5, 0.5, 0.5, 1.0},
		{-0.5, 0.5, 0.5, 1.0},
		{-0.5, 0.5, -0.5, 1.0},
		// bottom face
		{-0.5, -0.5, -0.5, 1.0},
		{0.5, -0.5, -0.5, 1.0},
		{0.5, -0.5, 0.5, 1.0},
		{0.5, -0.5, 0.5, 1.0},
		{-0.5, -0.5, 0.5, 1.0},
		{-0.5, -0.5, -0.5, 1.0},
	}

	idxVertices, indices := utils.GenerateIndexedVertices(cuboid)

	vertices := []float32{}
	for _, vertex := range idxVertices {
		v := c.transform.Mul4x1(vertex)
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

	var indexBuffer uint32
	gl.GenBuffers(1, &indexBuffer)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer)

	gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, nil)
}

func (c *Cuboid) UpdateTransform() {
	c.transform = math.MakeTransformationMatrix3D(c.Translation, c.Scalar, c.Rotation)
}

func (c *Cuboid) Translate(pos m.Vec3) {
	c.Translation = pos
	c.UpdateTransform()
}

func (c *Cuboid) Scale(scale m.Vec3) {
	c.Scalar = scale
	c.UpdateTransform()
}

func (c *Cuboid) Rotate(rotation m.Vec3) {
	c.Rotation = rotation
	c.UpdateTransform()
}
