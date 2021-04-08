package def

type Player struct {
	Name string `attr:"save"`
}

func init() {
	entities["Player"] = new(Player)
}
