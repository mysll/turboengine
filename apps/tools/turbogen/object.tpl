package {{.Pkg}}
import (
    . "turboengine/common/datatype"
    "turboengine/gameplay/object"
)

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
    {{$obj}}.InitOnce({{$obj}},{{len .Attrs}})
    {{range .Attrs}}{{$pri := "0"}}{{if eq .Save true}}{{$pri = (printf "%s|%s" $pri "object.OBJECT_SAVE")}}{{end}}
    {{if eq .Public true}}{{$pri = (printf "%s|%s" $pri "object.OBJECT_PUBLIC")}}{{end}}
    {{if eq .Private true}}{{$pri = (printf "%s|%s" $pri "object.OBJECT_PRIVATE")}}{{end}}
    {{if eq .Realtime true}}{{$pri = (printf "%s|%s" $pri "object.OBJECT_REALTIME")}}{{end}}
    {{if ne $pri "0"}}{{$obj}}.{{tolower .Name}}.SetFlag({{$pri}}){{end}}
    {{$obj}}.AddAttr({{$obj}}.{{tolower .Name}}){{end}}
    return {{$obj}}
}

{{range .Attrs}}
func ({{$obj}} *{{$.Name}}) {{.Name}}() {{alias .ArgType}} {
	return {{$obj}}.{{tolower .Name}}.Data()
}

func ({{$obj}} *{{$.Name}}) {{.Name}}Index() int {
    return {{$obj}}.{{tolower .Name}}.Index()
}

func ({{$obj}} *{{$.Name}}) Set{{.Name}}(v {{alias .ArgType}}) {
    {{if or (or (eq .Save true) (eq .Public true)) (eq .Private true)}}if {{$obj}}.{{tolower .Name}}.SetData(v) {
        {{if eq .Save true}}{{$obj}}.SetDirty()
        {{$obj}}.{{tolower .Name}}.SetFlag(object.OBJECT_DIRTY)
        {{end}}{{if eq .Public true}}{{$obj}}.SetPublicDirty()
        {{$obj}}.{{tolower .Name}}.SetFlag(object.OBJECT_PUBLIC_DIRTY)
        {{end}}{{if eq .Private true}}{{$obj}}.SetPrivateDirty()
        {{$obj}}.{{tolower .Name}}.SetFlag(object.OBJECT_PRIVATE_DIRTY){{end}} } {{else}}{{$obj}}.{{tolower .Name}}.SetData(v){{end}}
} 

func ({{$obj}} *{{$.Name}}) {{.Name}}Change(callback object.OnChange) {
    {{$obj}}.Change({{$obj}}.{{tolower .Name}}.Index(), callback)
}
{{end}}