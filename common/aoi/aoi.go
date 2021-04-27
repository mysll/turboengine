package aoi

import (
	. "turboengine/common/datatype"
)

type AOIMgr interface {
	Clear()
	GetIdsByRange(pos Vec3, ranges float32) []ObjectId
	Enter(obj ObjectId, pos Vec3, ranges float32)
	Leave(obj ObjectId, pos Vec3, ranges float32)
	Move(obj ObjectId, oldpos Vec3, dest Vec3, ranges float32) bool
}
