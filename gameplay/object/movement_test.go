package object

import (
	"testing"
)

func TestNewTransform(t *testing.T) {
	tr := NewTransform(nil)
	tr.SetRotation(Vec3{0, 0, 0})
	t.Log("forward:", tr.Forward())
	t.Log("up:", tr.Up())
	t.Log("right:", tr.Right())
	t.Log("rotation:", tr.rotation)
	t.Log("euler:", tr.EulerAngles())
	tr.LookAt(Vec3{1, 0, 1})
	t.Log("look at euler:", tr.EulerAngles())
	tr.SetRotation(Vec3{0, 90, 0})
	t.Log("rotation:", tr.rotation)
	t.Log("euler:", tr.EulerAngles())
	t.Log("forward:", tr.Forward())
}
