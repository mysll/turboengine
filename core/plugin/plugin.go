package plugin

import (
	"reflect"
	"turboengine/core/api"
)

var plugins = make(map[string]reflect.Type)

func Register(name string, p api.Plugin) {
	if _, dup := plugins[name]; dup {
		panic("plugin register twice " + name)
	}

	plugins[name] = reflect.TypeOf(p).Elem()
}

func NewPlugin(name string) api.Plugin {
	if p, ok := plugins[name]; ok {
		inst := reflect.New(p)
		p := inst.Interface().(api.Plugin)
		return p
	}
	return nil
}
