package level

import (
	_ "embed"
	"testing"
	"turboengine/common/datatype"
	"turboengine/gameplay/object"
)

//go:embed level.toml
var l string

type player struct {
	object.Object
}

func TestCreateFromFile(t *testing.T) {
	l := CreateFromData(l)
	ent := &player{}
	ent.SetId(1)
	ent.InitOnce(ent, 1)
	ent.SetFeature(object.FEATURES_ALL)
	ent.AOI().SetViewRange(100)
	l.AddEntity(ent)
	ent1 := &player{}
	ent1.SetId(2)
	ent1.InitOnce(ent, 1)
	ent1.SetFeature(object.FEATURES_ALL)
	ent1.AOI().SetViewRange(100)
	ent1.Movement().MoveTo(datatype.Vec3{50, 0, 50})
	l.AddEntity(ent1)
	t.Log(l.aoi.GetIdsByRange(datatype.Vec3{}, 100))
}
