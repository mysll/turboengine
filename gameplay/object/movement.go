package object

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Movement interface{}

// Transform 使用左手坐标系

type Transform struct {
	position mgl32.Vec3
	rotation mgl32.Quat // 四元数
	owner    GameObject
}

type EulerAngles struct {
	pitch float32 // x
	yaw   float32 // y
	roll  float32 // z
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
// see https://en.wikipedia.org/wiki/Conversion_between_quaternions_and_Euler_angles
func (t *Transform) Rotation() (pitch, yaw, roll float32) {
	return
}

func (t *Transform) SetRotation(pitch, yaw, roll float32) {
	t.rotation = mgl32.AnglesToQuat(mgl32.DegToRad(roll), mgl32.DegToRad(yaw), mgl32.DegToRad(pitch), mgl32.ZYX)
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
	t.rotation = mgl32.QuatLookAtV(t.position, mgl32.Vec3{x, y, z}, mgl32.Vec3{0, 1, 0})
}

func (t *Transform) MoveTo(x float32, y float32, z float32) {
	t.position = mgl32.Vec3{x, y, z}
}

func (t *Transform) Rotate(angle Vec3) {
}
