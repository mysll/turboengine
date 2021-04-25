package internal

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

func Clamp(value, min, max float64) float64 {
	return math.Max(min, math.Min(max, value))
}

func Euler(eulers mgl32.Vec3) mgl32.Quat {
	return mgl32.AnglesToQuat(mgl32.DegToRad(eulers.X()), mgl32.DegToRad(eulers.Y()), mgl32.DegToRad(eulers.Z()), mgl32.XYZ)
}

func QuatMulVec3(q mgl32.Quat, v mgl32.Vec3) mgl32.Vec3 {
	return q.Mul(mgl32.Quat{V: v}).Mul(q.Inverse()).V
}

func ToEuler(rotation mgl32.Quat) mgl32.Vec3 {
	var x, y, z float32

	r := rotation.Normalize()
	te := r.Mat4()
	m11 := float64(te[0])
	m12 := float64(te[4])
	m13 := float64(te[8])
	m22 := float64(te[5])
	m23 := float64(te[9])
	m32 := float64(te[6])
	m33 := float64(te[10])

	y = float32(math.Asin(Clamp(m13, -1, 1)))
	if math.Abs(m13) < 0.9999999 {
		x = float32(math.Atan2(-m23, m33))
		z = float32(math.Atan2(-m12, m11))
	} else {
		x = float32(math.Atan2(m32, m22))
		z = 0
	}

	return mgl32.Vec3{mgl32.RadToDeg(x), mgl32.RadToDeg(y), mgl32.RadToDeg(z)}
}

func LookAt(position, target mgl32.Vec3) mgl32.Quat {
	direction := target.Sub(position).Normalize()
	rotDir := mgl32.QuatBetweenVectors(mgl32.Vec3{0, 0, 1}, direction)
	return rotDir
}
