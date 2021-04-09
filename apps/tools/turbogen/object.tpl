package {{.Pkg}}
import "turboengine/gameplay/object"

type {{.Name}} struct {
	object.Object
	{{range .Attrs}}{{tolower .Name}} *{{getType .ArgType}}
    {{end}}
}
{{$obj := tolower $.Name}}
func New{{.Name}}() *{{.Name}} {
    {{$obj}}:=&{{.Name}}{
        {{range .Attrs}}{{tolower .Name}}: {{create .ArgType}}("{{.Name}}"),
        {{end}}
    }
    {{range .Attrs}}{{$obj}}.AddAttr({{$obj}}.{{tolower .Name}})
    {{end}}
    return {{$obj}}
}

{{range .Attrs}}
func ({{$obj}} *{{$.Name}}) {{.Name}}() {{.ArgType}} {
	return {{$obj}}.{{tolower .Name}}.Data()
}

func ({{$obj}} *{{$.Name}}) Set{{.Name}}(v {{.ArgType}}) {
	{{$obj}}.{{tolower .Name}}.SetData(v)
} 
{{end}}