package dao

var models = make(map[string]interface{})

func RegisterModel(name string, model interface{}) {
	if _, ok := models[name]; ok {
		panic("register model dup")
	}
	models[name] = model
}

func GetAllModel() map[string]interface{} {
	return models
}
