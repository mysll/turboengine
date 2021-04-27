package object

type Collider struct {
	owner GameObject
}

func NewCollider(owner GameObject) *Collider {
	return &Collider{
		owner: owner,
	}
}
