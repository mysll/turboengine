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
    {{$obj}}.New({{len .Attrs}})
    {{range .Attrs}}{{$pri := "0"}}{{if eq .Save true}}{{$pri = (printf "%s|%s" $pri "object.OBJECT_SAVE")}}{{end}}
    {{if eq .Public true}}{{$pri = (printf "%s|%s" $pri "object.OBJECT_PUBLIC")}}{{end}}
    {{if eq .Private true}}{{$pri = (printf "%s|%s" $pri "object.OBJECT_PRIVATE")}}{{end}}
    {{if eq .Realtime true}}{{$pri = (printf "%s|%s" $pri "object.OBJECT_REALTIME")}}{{end}}
    {{if ne $pri "0"}}{{$obj}}.{{tolower .Name}}.SetFlag({{$pri}}){{end}}
    {{$obj}}.AddAttr({{$obj}}.{{tolower .Name}}){{end}}
    {{$obj}}.Init()
    return {{$obj}}
}

{{range .Attrs}}
func ({{$obj}} *{{$.Name}}) {{.Name}}() {{.ArgType}} {
	return {{$obj}}.{{tolower .Name}}.Data()
}

func ({{$obj}} *{{$.Name}}) {{.Name}}Index() int {
    return {{$obj}}.{{tolower .Name}}.Index()
}

func ({{$obj}} *{{$.Name}}) Set{{.Name}}(v {{.ArgType}}) {
    {{if eq .Save true}}	    if {{$obj}}.{{tolower .Name}}.SetData(v) {
        {{$obj}}.SetDirty()
    } {{else}}
    {{$obj}}.{{tolower .Name}}.SetData(v){{end}}
} 

func ({{$obj}} *{{$.Name}}) {{.Name}}Change(callback object.OnChange) {
    {{$obj}}.Change({{$obj}}.{{tolower .Name}}.Index(), callback)
}
{{end}}