package object

import (
	"turboengine/gameplay/internal"

	"github.com/go-gl/mathgl/mgl32"
)

type Movement interface{}

// Transform 使用左手坐标系

type Transform struct {
	position mgl32.Vec3
	rotation mgl32.Quat // 四元数
	owner    GameObject
}

var (
	forward = mgl32.Vec3{0, 0, 1}
	up      = mgl32.Vec3{0, 1, 0}
	right   = mgl32.Vec3{1, 0, 0}
)

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

func (t *Transform) EulerAngles() Vec3 {
	return Vec3(internal.ToEuler(t.rotation))
}

func (t *Transform) SetRotation(eulers Vec3) {
	t.rotation = internal.Euler(mgl32.Vec3(eulers))
}

func (t *Transform) Forward() Vec3 {
	return Vec3(internal.QuatMulVec3(t.rotation, forward))
}

func (t *Transform) Up() Vec3 {
	return Vec3(internal.QuatMulVec3(t.rotation, up))
}

func (t *Transform) Right() Vec3 {
	return Vec3(internal.QuatMulVec3(t.rotation, right))
}

func (t *Transform) Translate(translation Vec3) {
	t.position = t.position.Add(mgl32.Vec3(translation))
}

func (t *Transform) LookAt(target Vec3) {
	t.rotation = internal.LookAt(t.position, mgl32.Vec3(target))
}

func (t *Transform) MoveTo(position Vec3) {
	t.position = mgl32.Vec3(position)
}

func (t *Transform) Rotate(eulers Vec3) {
	eulerRot := internal.Euler(mgl32.Vec3(eulers))
	t.rotation = t.rotation.Mul(t.rotation.Inverse().Mul(eulerRot).Mul(t.rotation))
}
