package math

import (
	m "github.com/go-gl/mathgl/mgl32"
)

func NormalizeVec2(v m.Vec2) m.Vec2 {
	return v.Mul(1.0 / v.Len())
}

func NormalizeVec3(v m.Vec3) m.Vec3 {
	return v.Mul(1.0 / v.Len())
}
