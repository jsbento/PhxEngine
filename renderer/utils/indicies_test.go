package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/davecgh/go-spew/spew"
	m "github.com/go-gl/mathgl/mgl32"
)

func TestGenerateIndices(t *testing.T) {
	vertices := []m.Vec4{
		{0, 0, 0, 1},
		{1, 0, 0, 1},
		{1, 1, 0, 1},
		{1, 1, 0, 1},
		{0, 1, 0, 1},
		{0, 0, 0, 1},
	}
	indices := GenerateIndices(vertices)
	spew.Dump(indices)
	assert.Equal(t, 6, len(indices))
}

func TestGenerateIndexedVertices(t *testing.T) {
	rawVertices := []m.Vec4{
		{0, 0, 0, 1},
		{1, 0, 0, 1},
		{1, 1, 0, 1},
		{1, 1, 0, 1},
		{0, 1, 0, 1},
		{0, 0, 0, 1},
	}
	exp := []m.Vec4{
		{0, 0, 0, 1},
		{1, 0, 0, 1},
		{1, 1, 0, 1},
		{0, 1, 0, 1},
	}
	vertices, indices := GenerateIndexedVertices(rawVertices)
	assert.Equal(t, 6, len(indices))

	for i, expVertex := range exp {
		assert.Equal(t, expVertex, vertices[i])
	}
}
