package utils

import (
	"math"

	m "github.com/go-gl/mathgl/mgl32"
)

type vec4Key struct {
	x uint32
	y uint32
	z uint32
	w uint32
}

func vec4KeyFromVec4(v m.Vec4) vec4Key {
	return vec4Key{
		x: math.Float32bits(v.X()),
		y: math.Float32bits(v.Y()),
		z: math.Float32bits(v.Z()),
		w: math.Float32bits(v.W()),
	}
}

func GenerateIndices(vertices []m.Vec4) []uint32 {
	indices := []uint32{}
	vertexToIndex := make(map[vec4Key]uint32)
	nextIdx := 0
	for _, vertex := range vertices {
		key := vec4KeyFromVec4(vertex)
		if idx, ok := vertexToIndex[key]; ok {
			indices = append(indices, idx)
		} else {
			vertexToIndex[key] = uint32(nextIdx)
			indices = append(indices, uint32(nextIdx))
			nextIdx++
		}
	}
	return indices
}

func GenerateIndexedVertices(rawVertices []m.Vec4) (vertices []m.Vec4, indices []uint32) {
	vertices, indices = []m.Vec4{}, []uint32{}
	vertexToIndex := make(map[vec4Key]uint32)
	nextIdx := 0

	for _, vertex := range rawVertices {
		key := vec4KeyFromVec4(vertex)
		if idx, ok := vertexToIndex[key]; ok {
			indices = append(indices, idx)
		} else {
			vertexToIndex[key] = uint32(nextIdx)
			indices = append(indices, uint32(nextIdx))
			vertices = append(vertices, vertex)
			nextIdx++
		}
	}

	return
}
