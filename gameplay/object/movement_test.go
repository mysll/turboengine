package object

import (
	"testing"
)

func TestNewTransform(t *testing.T) {
	tr := NewTransform(nil)
	tr.SetRotation(90, 90, 0)
	t.Log(tr.rotation)
	t.Log(tr.Rotation())
}
