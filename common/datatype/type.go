package datatype

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/mysll/toolkit"
	"math"
)

type ObjectId uint64

type Vec3 mgl32.Vec3

func V3(x float32, y float32, z float32) Vec3 {
	return Vec3{x, y, z}
}

func (v Vec3) X() float32 {
	return v[0]
}

func (v Vec3) Y() float32 {
	return v[1]
}

func (v Vec3) Z() float32 {
	return v[2]
}

func (v Vec3) Equal(rhs Vec3) bool {
	for i := 0; i < 3; i++ {
		if !toolkit.IsEqual32(v[i], rhs[i]) {
			return false
		}
	}
	return true
}

func (v Vec3) Mul(c float32) Vec3 {
	return Vec3{v[0] * c, v[1] * c, v[2] * c}
}

func (v Vec3) Add(v2 Vec3) Vec3 {
	return Vec3{v[0] + v2[0], v[1] + v2[1], v[2] + v2[2]}
}

func (v Vec3) Sub(v2 Vec3) Vec3 {
	return Vec3{v[0] - v2[0], v[1] - v2[1], v[2] - v2[2]}
}

func (v Vec3) Dot(v2 Vec3) float32 {
	return v[0]*v2[0] + v[1]*v2[1] + v[2]*v2[2]
}

func (v Vec3) Len() float32 {
	return float32(math.Sqrt(float64(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])))
}

func (v Vec3) LenSqr() float32 {
	return v[0]*v[0] + v[1]*v[1] + v[2]*v[2]
}

func (v Vec3) Normalize() Vec3 {
	l := 1.0 / v.Len()
	return Vec3{v[0] * l, v[1] * l, v[2] * l}
}

type Vec2 mgl32.Vec2

func V2(x float32, y float32) Vec2 {
	return Vec2{x, y}
}
