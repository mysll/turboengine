package internal

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

func Clamp(value, min, max float64) float64 {
	return math.Max(min, math.Min(max, value))
}

func Euler(eulers mgl32.Vec3) mgl32.Quat {
	return mgl32.AnglesToQuat(mgl32.DegToRad(eulers.Z()), mgl32.DegToRad(eulers.X()), mgl32.DegToRad(eulers.Y()), mgl32.ZXY)
}

func QuatMulVec3(q mgl32.Quat, v mgl32.Vec3) mgl32.Vec3 {
	return q.Mul(mgl32.Quat{V: v}).Mul(q.Inverse()).V
}

func toEulerYZX(rotation mgl32.Quat) mgl32.Vec3 {
	var x, y, z float32

	r := rotation.Normalize()
	te := r.Mat4()
	m11 := float64(te[0])
	//m12 := float64(te[4])
	m13 := float64(te[8])
	m21 := float64(te[1])
	m22 := float64(te[5])
	m23 := float64(te[9])
	m31 := float64(te[2])
	//m32 := float64(te[6])
	m33 := float64(te[10])

	// yzx
	z = float32(math.Asin(Clamp(m21, -1, 1)))
	if math.Abs(m21) < 0.9999999 {
		x = float32(math.Atan2(-m23, m22))
		y = float32(math.Atan2(-m31, m11))
	} else {
		x = float32(math.Atan2(m13, m33))
		y = 0
	}
	return mgl32.Vec3{mgl32.RadToDeg(x), mgl32.RadToDeg(y), mgl32.RadToDeg(z)}
}

func toEulerZXY(rotation mgl32.Quat) mgl32.Vec3 {
	var x, y, z float32

	r := rotation.Normalize()
	te := r.Mat4()
	m11 := float64(te[0])
	m12 := float64(te[4])
	//m13 := float64(te[8])
	m21 := float64(te[1])
	m22 := float64(te[5])
	//m23 := float64(te[9])
	m31 := float64(te[2])
	m32 := float64(te[6])
	m33 := float64(te[10])

	// yzx
	x = float32(math.Asin(Clamp(m32, -1, 1)))
	if math.Abs(m32) < 0.9999999 {
		y = float32(math.Atan2(-m31, m33))
		z = float32(math.Atan2(-m12, m22))
	} else {
		y = 0
		z = float32(math.Atan2(m21, m11))
	}
	return mgl32.Vec3{mgl32.RadToDeg(x), mgl32.RadToDeg(y), mgl32.RadToDeg(z)}
}

func ToEuler(rotation mgl32.Quat) mgl32.Vec3 {
	return toEulerZXY(rotation)
}

func LookAt(position, target mgl32.Vec3) mgl32.Quat {
	direction := target.Sub(position).Normalize()
	rotDir := mgl32.QuatBetweenVectors(mgl32.Vec3{0, 0, 1}, direction)
	return rotDir
}
