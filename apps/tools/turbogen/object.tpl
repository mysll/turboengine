package {{.Pkg}}
import "turboengine/gameplay/object"

type {{.Name}} struct {
	object.Object
	{{range .Attrs}}{{tolower .Name}} {{getType .ArgType}}
    {{end}}
}

{{range .Attrs}}
func (player *{{.Name}}) Name() string {
	return player.name.Data()
}

func (player *{{.Name}}) SetName(v string) {
	player.name.SetData(v)
} {{end}}