package object

type Movement interface{}

type Transform struct {
	X     float64
	Y     float64
	Z     float64
	owner interface{}
}

func NewTransform(owner interface{}) *Transform {
	return &Transform{
		owner: owner,
	}
}
