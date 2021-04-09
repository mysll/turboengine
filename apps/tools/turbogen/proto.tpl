package {{.Pkg}}

import(
	"{{.PkgPath}}"
	"fmt"
	"time"
	coreapi "turboengine/core/api"
	"turboengine/common/protocol"
	{{range $k, $v := .Imp}}"{{$k}}"{{end}}
)

type I{{.Name}}_RPC_Go_{{.Ver}} interface {
	{{range .Methods}}{{.Name}}({{range $k, $v := .ArgType}}{{if ne $k 0}},{{end}}{{$v}}{{end}})({{range $k, $v := .ReturnType}}{{if ne $k 0}},{{end}}{{$v}}{{end}})
	{{end}}
}

type {{.Name}}_RPC_Go_{{.Ver}} struct {
	handler I{{.Name}}_RPC_Go_{{.Ver}}
}

{{range .Methods}}
func (p *{{$.Name}}_RPC_Go_{{$.Ver}}) {{.Name}}(id uint16, data []byte) (ret *protocol.Message, err error) {
	ar := protocol.NewLoadArchive(data)
	{{range $k, $v := .ArgType}}
	var arg{{$k}} {{$v}}
	err = ar.Get(&arg{{$k}})
	if err != nil {
		return
	} {{end}}
	{{range $k, $v := .ReturnType}}{{if ne $v "error"}}reply{{$k}},{{end}}{{end}}err1 := p.handler.{{.Name}}({{range $k, $v := .ArgType}}{{if ne $k 0}},{{end}}arg{{$k}}{{end}})
	if err1 != nil {
		err = err1
		return
	}
	{{$l := len .ReturnType}}{{if gt $l 1}}//reply
	sr := protocol.NewAutoExtendArchive(128){{range $k, $v := .ReturnType}}{{if ne $v "error"}}
	err = sr.Put(reply{{$k}})
	if err != nil {
		return
	}{{end}}
	{{end}}
	ret = sr.Message()
	{{end}}
	return
}
{{end}}
func Set{{.Name}}Provider(svr coreapi.Service, prefix string, provider I{{.Name}}_RPC_Go_{{.Ver}}) error {
	m := new({{.Name}}_RPC_Go_{{.Ver}})
	m.handler = provider
	{{range .Methods}}
	if err := svr.Sub(fmt.Sprintf("%s%d:{{$.Name}}.{{.Name}}", prefix, svr.ID()), m.{{.Name}}); err != nil {
		return err
	}{{end}}
	return nil
}