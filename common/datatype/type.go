package datatype

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/mysll/toolkit"
)

type ObjectId uint64

type Vec3 mgl32.Vec3

var (
	Forward = Vec3{0, 0, 1}
	Back    = Vec3{0, 0, -1}
	Up      = Vec3{0, 1, 0}
	Down    = Vec3{0, -1, 0}
	Left    = Vec3{-1, 0, 0}
	Right   = Vec3{1, 0, 0}
	Zero    = Vec3{}
	One     = Vec3{1, 1, 1}
)

func V3(x float32, y float32, z float32) Vec3 {
	return Vec3{x, y, z}
}

func (v Vec3) X() float32 {
	return v[0]
}

func (v Vec3) Y() float32 {
	return v[1]
}

func (v Vec3) Z() float32 {
	return v[2]
}

func (v Vec3) Equal(rhs Vec3) bool {
	for i := 0; i < 3; i++ {
		if !toolkit.IsEqual32(v[i], rhs[i]) {
			return false
		}
	}
	return true
}

func (v Vec3) Cross(v2 Vec3) Vec3 {
	return Vec3{v[1]*v2[2] - v[2]*v2[1], v[2]*v2[0] - v[0]*v2[2], v[0]*v2[1] - v[1]*v2[0]}
}

func (v Vec3) Mul(c float32) Vec3 {
	return Vec3{v[0] * c, v[1] * c, v[2] * c}
}

func (v Vec3) Add(v2 Vec3) Vec3 {
	return Vec3{v[0] + v2[0], v[1] + v2[1], v[2] + v2[2]}
}

func (v Vec3) Sub(v2 Vec3) Vec3 {
	return Vec3{v[0] - v2[0], v[1] - v2[1], v[2] - v2[2]}
}

func (v Vec3) Dot(v2 Vec3) float32 {
	return v[0]*v2[0] + v[1]*v2[1] + v[2]*v2[2]
}

func (v Vec3) Len() float32 {
	return float32(math.Sqrt(float64(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])))
}

func (v Vec3) LenSqr() float32 {
	return v[0]*v[0] + v[1]*v[1] + v[2]*v[2]
}

func (v Vec3) Normalize() Vec3 {
	l := 1.0 / v.Len()
	return Vec3{v[0] * l, v[1] * l, v[2] * l}
}

// gorm
func (v *Vec3) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarsh Vec3", value))
	}
	result := Vec3{}
	err := json.Unmarshal(bytes, &result)
	if err != nil {
		return err
	}
	*v = result
	return nil
}

func (v Vec3) Value() (driver.Value, error) {
	return json.Marshal(v)
}

type Vec2 mgl32.Vec2

func V2(x float32, y float32) Vec2 {
	return Vec2{x, y}
}

func (v Vec2) Equal(rhs Vec2) bool {
	for i := 0; i < 2; i++ {
		if !toolkit.IsEqual32(v[i], rhs[i]) {
			return false
		}
	}
	return true
}

func (v Vec2) X() float32 {
	return v[0]
}

func (v Vec2) Y() float32 {
	return v[1]
}

func (v Vec2) Add(v2 Vec2) Vec2 {
	return Vec2{v[0] + v2[0], v[1] + v2[1]}
}

func (v Vec2) Sub(v2 Vec2) Vec2 {
	return Vec2{v[0] - v2[0], v[1] - v2[1]}
}

func (v Vec2) Mul(c float32) Vec2 {
	return Vec2{v[0] * c, v[1] * c}
}

func (v Vec2) Dot(v2 Vec2) float32 {
	return v[0]*v2[0] + v[1]*v2[1]
}

func (v Vec2) Len() float32 {
	return float32(math.Hypot(float64(v[0]), float64(v[1])))
}

func (v Vec2) LenSqr() float32 {
	return v[0]*v[0] + v[1]*v[1]
}

func (v Vec2) Normalize() Vec2 {
	l := 1.0 / v.Len()
	return Vec2{v[0] * l, v[1] * l}
}

// gorm
func (v *Vec2) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarsh Vec3", value))
	}
	result := Vec2{}
	err := json.Unmarshal(bytes, &result)
	if err != nil {
		return err
	}
	*v = result
	return nil
}

func (v Vec2) Value() (driver.Value, error) {
	return json.Marshal(v)
}
