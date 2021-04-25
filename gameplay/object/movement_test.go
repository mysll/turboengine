package object

import (
	"testing"
)

func TestNewTransform(t *testing.T) {
	tr := NewTransform(nil)
	tr.SetRotation(Vec3{90, 0, 0})
	t.Log(tr.rotation)
}
