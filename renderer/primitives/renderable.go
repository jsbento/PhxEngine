package primitives

import (
	m "github.com/go-gl/mathgl/mgl32"
)

type Renderable2D interface {
	Draw()
	Translate(pos m.Vec2)
	Scale(scale m.Vec2)
	Rotate(angle float32)
	UpdateTransform()
}

type Renderable3D interface {
	Draw()
	Translate(pos m.Vec3)
	Scale(scale m.Vec3)
	Rotate(rotation m.Vec3)
	UpdateTransform()
}
