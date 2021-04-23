package object

import "github.com/go-gl/mathgl/mgl32"

type Movement interface{}

// Transform 使用左手坐标系

type Transform struct {
	position mgl32.Vec3
	rotation mgl32.Quat // 四元数
	owner    GameObject
}

type EulerAngles struct {
	roll  float32
	pitch float32
	yaw   float32
}

func NewTransform(owner GameObject) *Transform {
	return &Transform{
		owner:    owner,
		position: mgl32.Vec3{},
		rotation: mgl32.QuatIdent(),
	}
}

func (t *Transform) Position() Vec3 {
	return Vec3(t.position)
}

// Rotation get euler angle (roll pitch yaw)
func (t *Transform) Rotation() (roll, pitch, yaw float32) {
	return
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

func (t *Transform) MoveTo(x float32, y float32, z float32) {
	t.position = mgl32.Vec3{x, y, z}
}

func (t *Transform) Rotate(angle Vec3) {

}
