package internal

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

const order = mgl32.ZXY

func Clamp(value, min, max float64) float64 {
	return math.Max(min, math.Min(max, value))
}

func Euler(eulers mgl32.Vec3) mgl32.Quat {
	switch order {
	case mgl32.ZXY:
		return mgl32.AnglesToQuat(mgl32.DegToRad(eulers.Z()), mgl32.DegToRad(eulers.X()), mgl32.DegToRad(eulers.Y()), mgl32.ZXY)
	case mgl32.YZX:
		return mgl32.AnglesToQuat(mgl32.DegToRad(eulers.Y()), mgl32.DegToRad(eulers.Z()), mgl32.DegToRad(eulers.X()), mgl32.YZX)
	}
	panic("unsupport order")
}

func QuatMulVec3(q mgl32.Quat, v mgl32.Vec3) mgl32.Vec3 {
	return q.Mul(mgl32.Quat{V: v}).Mul(q.Inverse()).V
}

func toEulerYZX(rotation mgl32.Quat) mgl32.Vec3 {
	var x, y, z, w float32

	r := rotation.Normalize()
	w, x, y, z = r.W, r.V[0], r.V[1], r.V[2]
	m11 := float64(1 - 2*y*y - 2*z*z)
	m13 := float64(2*x*z + 2*w*y)
	m21 := float64(2*x*y + 2*w*z)
	m22 := float64(1 - 2*x*x - 2*z*z)
	m23 := float64(2*y*z - 2*w*x)
	m31 := float64(2*x*z - 2*w*y)
	m33 := float64(1 - 2*x*x - 2*y*y)

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
	var x, y, z, w float32

	r := rotation.Normalize()
	w, x, y, z = r.W, r.V[0], r.V[1], r.V[2]
	m11 := float64(1 - 2*y*y - 2*z*z)
	m12 := float64(2*x*y - 2*w*z)
	m21 := float64(2*x*y + 2*w*z)
	m22 := float64(1 - 2*x*x - 2*z*z)
	m31 := float64(2*x*z - 2*w*y)
	m32 := float64(2*y*z + 2*w*x)
	m33 := float64(1 - 2*x*x - 2*y*y)

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
	switch order {
	case mgl32.ZXY:
		return toEulerZXY(rotation)
	case mgl32.YZX:
		return toEulerYZX(rotation)
	}
	panic("unsupport order")
}

func LookAt(position, target mgl32.Vec3) mgl32.Quat {
	direction := target.Sub(position).Normalize()
	return mgl32.QuatBetweenVectors(mgl32.Vec3{0, 0, 1}, direction)
}
