package {{.Pkg}}
import "turboengine/gameplay/object"

type {{.Name}} struct {
	object.Object
	{{range .Attrs}}{{tolower .Name}} {{getType .ArgType}}
    {{end}}
}

{{$obj := tolower $.Name}}
{{range .Attrs}}
func ({{$obj}} *{{.Name}}) {{.Name}}() {{.ArgType}} {
	return {{$obj}}.{{tolower .Name}}.Data()
}

func ({{$obj}} *{{.Name}}) Set{{.Name}}(v {{.ArgType}}) {
	{{$obj}}.{{tolower .Name}}.SetData(v)
} {{end}}