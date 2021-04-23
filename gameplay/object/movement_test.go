package object

import (
	"testing"
)

func TestNewTransform(t *testing.T) {
	tr := NewTransform(nil)
	t.Log(tr.Position())
}
