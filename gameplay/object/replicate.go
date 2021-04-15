package object

import "io"

type Replicate interface {
	WriteAll(io.Writer) (int, error)
	WritePublic(io.Writer) (int, error)
	WritePrivate(io.Writer) (int, error)
}

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

func (r *Replication) WriteAll(stream io.Writer) (int, error) {
	return 0, nil
}

func (r *Replication) WritePublic(stream io.Writer) (int, error) {
	return 0, nil
}

func (r *Replication) WritePrivate(stream io.Writer) (int, error) {
	return 0, nil
}
