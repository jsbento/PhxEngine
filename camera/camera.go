package camera

import (
	m "github.com/go-gl/mathgl/mgl32"
)

type CameraType string

const (
	CameraTypeOrthographic CameraType = "Orthographic"
	CameraTypePerspective  CameraType = "Perspective"
)

type Camera struct {
	Type CameraType

	PerspectiveFOV  float32
	PerspectiveNear float32
	PerspectiveFar  float32

	OrthographicSize float32
	OrthographicNear float32
	OrthographicFar  float32

	AspectRatio float32

	ViewMatrix     m.Mat4
	Projection     m.Mat4
	ViewProjection m.Mat4
}

func NewOrthographicCamera(size, near, far float32) (c *Camera) {
	c = &Camera{
		Type:             CameraTypeOrthographic,
		OrthographicSize: size,
		OrthographicNear: near,
		OrthographicFar:  far,
	}
	c.RecalculateProjection()
	return
}

func NewPerspectiveCamera(fov, near, far float32) (c *Camera) {
	c = &Camera{
		Type:            CameraTypePerspective,
		PerspectiveFOV:  fov,
		PerspectiveNear: near,
		PerspectiveFar:  far,
	}
	c.RecalculateProjection()
	return
}

func (c *Camera) SetViewportSize(width, height float32) {
	c.AspectRatio = width / height
	c.RecalculateProjection()
}

func (c *Camera) RecalculateProjection() {
	switch c.Type {
	case CameraTypeOrthographic:
		orthoLeft := -c.OrthographicSize * c.AspectRatio * 0.5
		orthoRight := c.OrthographicSize * c.AspectRatio * 0.5
		orthoBottom := -c.OrthographicSize * 0.5
		orthoTop := c.OrthographicSize * 0.5

		c.Projection = m.Ortho(
			orthoLeft,
			orthoRight,
			orthoBottom,
			orthoTop,
			c.OrthographicNear,
			c.OrthographicFar,
		)
	case CameraTypePerspective:
		c.Projection = m.Perspective(
			c.PerspectiveFOV,
			c.AspectRatio,
			c.PerspectiveNear,
			c.PerspectiveFar,
		)
	}
}
