package {{.Pkg}}

type {{.Name}} struct {
	Name  string `attr:"save"`
}

func init() {
	entities["{{.Name}}"] = new({{.Name}})
}