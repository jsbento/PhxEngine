package primitives

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	m "github.com/go-gl/mathgl/mgl32"
	"github.com/jsbento/PhxEngine/math"
	"github.com/jsbento/PhxEngine/renderer/utils"
)

type Cuboid struct {
	renderId    uint32
	vao         uint32
	indexBuffer uint32
	indexCount  int32
	baseVerts   []m.Vec4
	indexData   []uint32
	vertexData  []float32
	buffersInit bool
	dirty       bool
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
	c.baseVerts = idxVertices
	c.indexData = indices
	c.indexCount = int32(len(indices))

	c.UpdateTransform()
	return
}

func (c *Cuboid) Draw() {
	c.ensureBuffers()
	c.updateVertexData()
	gl.BindVertexArray(c.vao)
	gl.DrawElements(gl.TRIANGLES, c.indexCount, gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)
}

func (c *Cuboid) UpdateTransform() {
	c.transform = math.MakeTransformationMatrix3D(c.Translation, c.Scalar, c.Rotation)
	c.dirty = true
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

func (c *Cuboid) ensureBuffers() {
	if c.buffersInit {
		return
	}
	gl.GenVertexArrays(1, &c.vao)
	gl.GenBuffers(1, &c.renderId)
	gl.GenBuffers(1, &c.indexBuffer)

	gl.BindVertexArray(c.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, c.renderId)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, c.indexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(c.indexData)*4, gl.Ptr(c.indexData), gl.STATIC_DRAW)

	gl.BindVertexArray(0)
	c.buffersInit = true
}

func (c *Cuboid) updateVertexData() {
	if !c.dirty {
		return
	}
	required := len(c.baseVerts) * 3
	if cap(c.vertexData) < required {
		c.vertexData = make([]float32, 0, required)
	} else {
		c.vertexData = c.vertexData[:0]
	}
	for _, vertex := range c.baseVerts {
		v := c.transform.Mul4x1(vertex)
		c.vertexData = append(c.vertexData, v.X(), v.Y(), v.Z())
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, c.renderId)
	gl.BufferData(gl.ARRAY_BUFFER, len(c.vertexData)*4, gl.Ptr(c.vertexData), gl.DYNAMIC_DRAW)
	c.dirty = false
}

func (c *Cuboid) Destroy() {
	if !c.buffersInit {
		return
	}
	gl.DeleteBuffers(1, &c.renderId)
	gl.DeleteBuffers(1, &c.indexBuffer)
	gl.DeleteVertexArrays(1, &c.vao)
	c.renderId = 0
	c.indexBuffer = 0
	c.vao = 0
	c.buffersInit = false
}
