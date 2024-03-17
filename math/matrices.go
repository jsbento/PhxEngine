package math

import (
	m "github.com/go-gl/mathgl/mgl32"
)

func MakeScaleMatrix(scale m.Vec3) m.Mat4 {
	return m.Scale3D(scale.X(), scale.Y(), scale.Z())
}

// angle in radians
func MakeRotationMatrix(angle float32, axis m.Vec3) m.Mat4 {
	return m.HomogRotate3D(angle, axis)
}

func MakeRotationMatrix3D(rotation m.Vec3) m.Mat4 {
	return m.HomogRotate3D(rotation.X(), m.Vec3{1, 0, 0}).
		Mul4(m.HomogRotate3D(rotation.Y(), m.Vec3{0, 1, 0})).
		Mul4(m.HomogRotate3D(rotation.Z(), m.Vec3{0, 0, 1}))
}

func MakeTranslationMatrix(translation m.Vec3) m.Mat4 {
	return m.Translate3D(translation.X(), translation.Y(), translation.Z())
}

func MakeTransformationMatrix(scale m.Vec3, angle float32, axis m.Vec3, translation m.Vec3) m.Mat4 {
	return MakeTranslationMatrix(
		translation,
	).Mul4(MakeRotationMatrix(angle, axis)).
		Mul4(MakeScaleMatrix(scale))
}

func MakeTransformationMatrix3D(translation m.Vec3, scale m.Vec3, rotation m.Vec3) m.Mat4 {
	return MakeTranslationMatrix(translation).
		Mul4(MakeRotationMatrix3D(rotation)).
		Mul4(MakeScaleMatrix(scale))
}
