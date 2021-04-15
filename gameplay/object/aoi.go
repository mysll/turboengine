package object

type AOI interface {
	Snapshot() *ViewChange
	Clear()
}

type ViewChange struct {
	News []GameObject
	Del  []GameObject
}

type View struct {
	around []GameObject
	owner  GameObject
}

func NewView(owner GameObject) *View {
	return &View{
		owner:  owner,
		around: make([]GameObject, 30),
	}
}

func (v *View) Clear() {

}

// 捕获当前周围的环境
func (v *View) Snapshot() *ViewChange {
	return nil
}
