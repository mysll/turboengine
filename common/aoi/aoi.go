package aoi

import "turboengine/gameplay/object"

type AOIMgr interface {
	Clear()
	GetIdsByPos(pos object.Vec3, ranges int) []object.ObjectId
	AddObject(pos object.Vec3, obj object.ObjectId) bool
	RemoveObject(pos object.Vec3, obj object.ObjectId) bool
	UpdateObject(obj object.ObjectId, oldpos object.Vec3, newpos object.Vec3) error
	GetWatchers(pos object.Vec3) []object.ObjectId
	AddWatcher(watcher object.ObjectId, pos object.Vec3, ranges int) bool
	RemoveWatcher(watcher object.ObjectId, pos object.Vec3, ranges int) bool
	UpdateWatcher(watcher object.ObjectId, oldPos object.Vec3, newPos object.Vec3, oldRange, newRange int) bool
}
