package def

import (
	"reflect"
	"testing"
	"turboengine/apps/tools/turbogen"
)

func TestCreate(t *testing.T) {
	for _, v := range entities {
		typ := reflect.TypeOf(v)
		turbogen.ObjectWrap(v, typ.Elem().PkgPath(), "entity", "../entity")
	}
}
