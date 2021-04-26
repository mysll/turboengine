package aoi

import "turboengine/gameplay/object"

type AOIMgr interface {
	Clear()
	GetIdsByRange(pos object.Vec3, ranges int) []object.ObjectId
	Enter(obj object.ObjectId, pos object.Vec3, ranges int)
	Level(obj object.ObjectId, pos object.Vec3, ranges int)
	Move(obj object.ObjectId, oldpos object.Vec3, dest object.Vec3, ranges int) bool
}
