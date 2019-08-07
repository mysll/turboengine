package proto

import (
	"testing"
	"turboengine/apps/tools/turbogen"
)

func TestCreate(t *testing.T) {
	for _, v := range reg {
		turbogen.Generate(v, "rpc", "../rpc")
	}
}
