package object

import "github.com/go-gl/mathgl/mgl32"

type Movement interface{}

type Transform struct {
	pos    mgl32.Vec3
	orient float32 // 0~360
	owner  GameObject
}

func NewTransform(owner GameObject) *Transform {
	return &Transform{
		owner: owner,
	}
}

func (t *Transform) MoveTo(x float32, y float32, z float32, orient float32) {
	t.pos = mgl32.Vec3{x, y, z}
	t.orient = orient
}

func (t *Transform) Forward() Vec3 {
	return Vec3{}
}

func (t *Transform) Up() Vec3 {
	return Vec3{}
}

func (t *Transform) Right() Vec3 {
	return Vec3{}
}

func (t *Transform) Translate(x float32, y float32, z float32) {

}

func (t *Transform) LookAt(x float32, y float32, z float32) {

}
