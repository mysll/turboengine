package def

import "turboengine/common/datatype"

type Player struct {
	Name  string `attr:"save,public" orm:"uniqueIndex;size:64"`
	Sex   int32  `attr:"save,private" orm:"size:1"`
	Hp    int64  `attr:"public,realtime"`
	Dmg   float32
	Test  float64
	Pos   datatype.Vec3 `attr:"save,public,private,realtime"  orm:"type:varbinary(64)"`
	Bytes []byte        `attr:"save" orm:"size:64"`
}

func init() {
	entities["Player"] = new(Player)
}
