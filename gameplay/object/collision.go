package object

import (
	. "turboengine/common/datatype"
)

type Shape interface {
	Init(args ...interface{})
	Type() int
	SetCenter(pos Vec3)
}

type Collider interface {
	Collide(other Collider) bool
	Shape() Shape
}

type Collision struct {
	owner GameObject
	shape Shape
}

func NewCollision(owner GameObject) *Collision {
	return &Collision{
		owner: owner,
	}
}

func (c *Collision) Shape() Shape {
	if c.shape != nil {
		c.shape.SetCenter(c.owner.Movement().Position())
	}
	return c.shape
}

func (c *Collision) Collide(other Collider) bool {
	return false
}
