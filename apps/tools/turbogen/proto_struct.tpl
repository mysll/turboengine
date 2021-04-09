package {{.Pkg}}

type {{.Name}} struct {
	Ver   string {{.Tag}}
	XXX   interface{}
	// custom method begin

	// custom method end
}

func init() {
	reg["{{.Name}}"] = new({{.Name}})
}