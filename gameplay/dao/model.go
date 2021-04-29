package dao

var models = make(map[string]Persistent)

type Persistent interface {
	DBId() uint64
}

func RegisterModel(name string, model Persistent) {
	if _, ok := models[name]; ok {
		panic("register model dup")
	}
	models[name] = model
}

func GetAllModel() map[string]Persistent {
	return models
}
