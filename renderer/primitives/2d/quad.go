package primitives

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	m "github.com/go-gl/mathgl/mgl32"
	"github.com/jsbento/PhxEngine/math"
	"github.com/jsbento/PhxEngine/renderer/utils"
)

type Quad struct {
	renderId    uint32
	transform   m.Mat4
	Translation m.Vec2
	Scalar      m.Vec2
	Rotation    float32
}

func NewQuad(
	translation m.Vec2,
	scale m.Vec2,
	rotation float32,
) (q *Quad) {
	q = &Quad{
		Translation: translation,
		Scalar:      scale,
		Rotation:    rotation,
	}
	q.UpdateTransform()
	return
}

func (q *Quad) Draw() {
	square := []m.Vec4{
		{-0.5, -0.5, 0.0, 1.0},
		{0.5, -0.5, 0.0, 1.0},
		{0.5, 0.5, 0.0, 1.0},
		{0.5, 0.5, 0.0, 1.0},
		{-0.5, 0.5, 0.0, 1.0},
		{-0.5, -0.5, 0.0, 1.0},
	}

	idxVertices, indices := utils.GenerateIndexedVertices(square)

	vertices := []float32{}
	for _, vertex := range idxVertices {
		v := q.transform.Mul4x1(vertex)
		vertices = append(vertices, v.X(), v.Y(), v.Z())
	}

	gl.GenBuffers(1, &q.renderId)
	gl.BindBuffer(gl.ARRAY_BUFFER, q.renderId)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, q.renderId)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	var indexBuffer uint32
	gl.GenBuffers(1, &indexBuffer)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer)

	gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, nil)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (q *Quad) Translate(pos m.Vec2) {
	q.Translation = pos
	q.UpdateTransform()
}

func (q *Quad) Scale(scale m.Vec2) {
	q.Scalar = scale
	q.UpdateTransform()
}

func (q *Quad) Rotate(angle float32) {
	q.Rotation = angle
	q.UpdateTransform()
}

func (q *Quad) UpdateTransform() {
	q.transform = math.MakeTransformationMatrix(
		q.Scalar.Vec3(1.0),
		q.Rotation,
		m.Vec3{0.0, 0.0, 1.0},
		q.Translation.Vec3(0.0),
	)
}
