package internal

import "github.com/go-gl/mathgl/mgl32"

func Euler(eulers mgl32.Vec3) mgl32.Quat {
	return mgl32.AnglesToQuat(mgl32.DegToRad(eulers.X()), mgl32.DegToRad(eulers.Y()), mgl32.DegToRad(eulers.Z()), mgl32.XYZ)
}

func QuatMulVec3(q mgl32.Quat, v mgl32.Vec3) mgl32.Vec3 {
	return q.Mul(mgl32.Quat{V: v}).Mul(q.Inverse()).V
}

func ToEuler(q mgl32.Quat) mgl32.Vec3 {
	return mgl32.Vec3{}
}
