package def

import "turboengine/common/datatype"

type Player struct {
	Name  string `attr:"save,public"`
	Sex   int32  `attr:"save,private"`
	Hp    int64  `attr:"public,realtime"`
	Dmg   float32
	Test  float64
	Pos   datatype.Vec3 `attr:"public,private,realtime"`
	Bytes []byte
}

func init() {
	entities["Player"] = new(Player)
}
