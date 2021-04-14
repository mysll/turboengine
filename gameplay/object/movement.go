package object

type Movement interface{}

type Transform struct {
	X     float64
	Y     float64
	Z     float64
	owner GameObject
}

func NewTransform(owner GameObject) *Transform {
	return &Transform{
		owner: owner,
	}
}
