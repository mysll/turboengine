package internal

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestQuatMulVec3(t *testing.T) {
	q := Euler(mgl32.Vec3{90, 90, 90})
	t.Log(q)
	v := QuatMulVec3(q, mgl32.Vec3{0, 0, 10})
	t.Log(v)
	euler := ToEuler(q)
	t.Log(euler)
}
