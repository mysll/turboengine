package proto

import (
	"reflect"
	"testing"
	"turboengine/apps/tools/turbogen"
)

func TestCreate(t *testing.T) {
	for _, v := range reg {
		typ := reflect.TypeOf(v)
		turbogen.Generate(v, typ.Elem().PkgPath(), "rpc", "../rpc", "Gate")
	}
}
