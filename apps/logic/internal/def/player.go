package def

import "turboengine/gameplay/object"

type Player struct {
	Name string `attr:"save,public"`
	Sex  int32  `attr:"save,private"`
	Hp   int64  `attr:"public,realtime"`
	Dmg  float32
	Test float64
	Pos  object.Vec3 `attr:"public,private,realtime"`
}

func init() {
	entities["Player"] = new(Player)
}
