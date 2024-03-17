package utils

import (
	"fmt"

	m "github.com/go-gl/mathgl/mgl32"
)

func GenerateIndices(vertices []m.Vec4) []uint32 {
	indices := []uint32{}
	vertexToIndex := make(map[string]uint32)
	nextIdx := 0
	for _, vertex := range vertices {
		key := stringifyVec4(vertex)
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
	vertexToIndex := make(map[string]uint32)
	nextIdx := 0

	for _, vertex := range rawVertices {
		key := stringifyVec4(vertex)
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

func stringifyVec4(v m.Vec4) string {
	return fmt.Sprintf("{%f, %f, %f, %f}", v.X(), v.Y(), v.Z(), v.W())
}
