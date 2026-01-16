package primitives

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	m "github.com/go-gl/mathgl/mgl32"
	"github.com/jsbento/PhxEngine/math"
	"github.com/jsbento/PhxEngine/renderer/utils"
)

type Quad struct {
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
	square := []m.Vec4{
		{-0.5, -0.5, 0.0, 1.0},
		{0.5, -0.5, 0.0, 1.0},
		{0.5, 0.5, 0.0, 1.0},
		{0.5, 0.5, 0.0, 1.0},
		{-0.5, 0.5, 0.0, 1.0},
		{-0.5, -0.5, 0.0, 1.0},
	}

	idxVertices, indices := utils.GenerateIndexedVertices(square)
	q.baseVerts = idxVertices
	q.indexData = indices
	q.indexCount = int32(len(indices))

	q.UpdateTransform()
	return
}

func (q *Quad) Draw() {
	q.ensureBuffers()
	q.updateVertexData()
	gl.BindVertexArray(q.vao)
	gl.DrawElements(gl.TRIANGLES, q.indexCount, gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)
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
	q.dirty = true
}

func (q *Quad) ensureBuffers() {
	if q.buffersInit {
		return
	}
	gl.GenVertexArrays(1, &q.vao)
	gl.GenBuffers(1, &q.renderId)
	gl.GenBuffers(1, &q.indexBuffer)

	gl.BindVertexArray(q.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, q.renderId)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, q.indexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(q.indexData)*4, gl.Ptr(q.indexData), gl.STATIC_DRAW)

	gl.BindVertexArray(0)
	q.buffersInit = true
}

func (q *Quad) updateVertexData() {
	if !q.dirty {
		return
	}
	required := len(q.baseVerts) * 3
	if cap(q.vertexData) < required {
		q.vertexData = make([]float32, 0, required)
	} else {
		q.vertexData = q.vertexData[:0]
	}
	for _, vertex := range q.baseVerts {
		v := q.transform.Mul4x1(vertex)
		q.vertexData = append(q.vertexData, v.X(), v.Y(), v.Z())
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, q.renderId)
	gl.BufferData(gl.ARRAY_BUFFER, len(q.vertexData)*4, gl.Ptr(q.vertexData), gl.DYNAMIC_DRAW)
	q.dirty = false
}

func (q *Quad) Destroy() {
	if !q.buffersInit {
		return
	}
	gl.DeleteBuffers(1, &q.renderId)
	gl.DeleteBuffers(1, &q.indexBuffer)
	gl.DeleteVertexArrays(1, &q.vao)
	q.renderId = 0
	q.indexBuffer = 0
	q.vao = 0
	q.buffersInit = false
}
