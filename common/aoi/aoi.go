package aoi

import (
	. "turboengine/common/datatype"
)

type AOIMgr interface {
	Clear()
	GetIdsByRange(pos Vec3, ranges int) []ObjectId
	Enter(obj ObjectId, pos Vec3, ranges int)
	Leave(obj ObjectId, pos Vec3, ranges int)
	Move(obj ObjectId, oldpos Vec3, dest Vec3, ranges int) bool
}
