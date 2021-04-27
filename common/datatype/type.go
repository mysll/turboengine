package datatype

import "github.com/go-gl/mathgl/mgl32"

type ObjectId uint64

type Vec3 mgl32.Vec3

func V3(x float32, y float32, z float32) Vec3 {
	return [3]float32{x, y, z}
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

type Vec2 mgl32.Vec2

func V2(x float32, y float32) Vec2 {
	return [2]float32{x, y}
}

func (v Vec2) X() float32 {
	return v[0]
}

func (v Vec2) Y() float32 {
	return v[1]
}
