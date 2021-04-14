package object

type Replicate interface{}

type Replication struct {
	owner interface{}
}

func NewReplication(owner interface{}) *Replication {
	return &Replication{
		owner: owner,
	}
}
