package aoi

import "turboengine/gameplay/object"

type AOIMgr interface {
	Clear()
	GetIdsByPos(pos object.Vec3, ranges int) []object.ObjectId
	GetIdsByType(pos object.Vec3, ranges int, types []int) []object.ObjectId
	AddObject(pos object.Vec3, obj object.ObjectId, typ int) bool
	RemoveObject(pos object.Vec3, obj object.ObjectId, typ int) bool
	UpdateObject(obj object.ObjectId, typ int, oldpos object.Vec3, newpos object.Vec3) error
	GetWatchers(pos object.Vec3, types []int) []object.ObjectId
	AddWatcher(watcher object.ObjectId, typ int, pos object.Vec3, ranges int) bool
	RemoveWatcher(watcher object.ObjectId, typ int, pos object.Vec3, ranges int) bool
	UpdateWatcher(watcher object.ObjectId, typ int, oldPos object.Vec3, newPos object.Vec3, oldRange, newRange int) bool
}
