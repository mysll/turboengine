package datatype

import "github.com/go-gl/mathgl/mgl32"

type ObjectId uint64

type Vec3 struct {
	mgl32.Vec3
}

func V3(x float32, y float32, z float32) Vec3 {
	return Vec3{
		Vec3: mgl32.Vec3{x, y, z},
	}
}

type Vec2 struct {
	mgl32.Vec2
}

func V2(x float32, y float32) Vec2 {
	return Vec2{
		Vec2: mgl32.Vec2{x, y},
	}
}
