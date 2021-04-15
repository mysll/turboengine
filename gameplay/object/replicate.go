package object

type Replicate interface{}

type Replication struct {
	owner GameObject
}

func NewReplication(owner GameObject) *Replication {
	return &Replication{
		owner: owner,
	}
}

func (r *Replication) change(index int, val interface{}) {

}
