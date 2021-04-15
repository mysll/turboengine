package object

type Movement interface{}

type Transform struct {
	X      float32
	Y      float32
	Z      float32
	Orient float32
	owner  GameObject
}

func NewTransform(owner GameObject) *Transform {
	return &Transform{
		owner: owner,
	}
}

func (t *Transform) MoveTo(x float32, y float32, z float32, orient float32) {
	t.X, t.Y, t.Z, t.Orient = x, y, z, orient
}
