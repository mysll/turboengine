package object

import (
	. "turboengine/common/datatype"
)

type Shape interface {
	Init(args ...interface{})
	Type() int
	SetCenter(pos Vec3)
	Center() Vec3
	Collider(Shape) bool
}

const (
	SHAPE_BOX = 1 + iota
	SHAPE_CYLINDER
	SHAPE_CAPSULE
	SHAPE_CIRCLE_2D
)

type shape struct {
	center Vec3
}

func (s *shape) SetCenter(c Vec3) {
	s.center = c
}

func (s *shape) Center() Vec3 {
	return s.center
}

type Circle2DShape struct {
	shape
	radius float32
}

func (s *Circle2DShape) Init(args ...interface{}) {
	switch v := args[0].(type) {
	case float32:
		s.radius = v
	}
}

func (s *Circle2DShape) Type() int {
	return SHAPE_CIRCLE_2D
}

func (s *Circle2DShape) SetCenter(c Vec3) {
	s.center = Vec3{c.X(), 0, c.Z()}
}

func (s *Circle2DShape) Collider(other Shape) bool {
	switch rhs := other.(type) {
	case *Circle2DShape:
		return s.center.Sub(rhs.center).Len() < s.radius+rhs.radius
		//TODO check other shape
	default:
		return false
	}
}

//TODO boxshape
type BoxShape struct {
	shape
}

//TODO CylinderShape
type CylinderShape struct {
	shape
}

// TODO other

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
	lhs := c.Shape()
	rhs := other.Shape()
	return lhs.Collider(rhs)
}
