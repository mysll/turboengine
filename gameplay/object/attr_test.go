package object

import "testing"

func TestAttrHolder_All(t *testing.T) {
	attr := AttrHolder{}
	attr.SetFlag(1)
	if attr.Flag() != 1 {
		t.Fatal("not match")
	}
	attr.SetFlag(2)
	attr.ClearFlag(1)
	if attr.Flag() != 2 {
		t.Fatal("not match")
	}
}
