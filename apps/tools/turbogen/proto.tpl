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

// client
type {{$.Name}}_RPC_Go_{{$.Ver}}_Client struct {
	svr coreapi.Service
	prefix string
	dest    protocol.Mailbox
	timeout time.Duration
	selector coreapi.Selector
}

func (m *{{$.Name}}_RPC_Go_{{$.Ver}}_Client) Redirect(dest protocol.Mailbox) {
	m.dest = dest
}

func (m *{{$.Name}}_RPC_Go_{{$.Ver}}_Client) SetSelector(selector coreapi.Selector) {
	m.selector = selector
}

{{range .Methods}}
// {{.Name}}
func (m *{{$.Name}}_RPC_Go_{{$.Ver}}_Client) {{.Name}}({{range $k, $v := .ArgType}}{{if ne $k 0}},{{end}}arg{{$k}} {{$v}}{{end}}) ({{range $k, $v := .ReturnType}}{{if ne $k 0}},{{end}}{{if eq $v "error"}}err{{else}}reply{{$k}}{{end}} {{$v}}{{end}}) {
	sr := protocol.NewAutoExtendArchive(128)
	{{range $k, $v := .ArgType}}err = sr.Put(arg{{$k}})
	if err != nil {
		return
	}
	{{end}}
	msg := sr.Message()
	remote := m.dest
	if remote.IsNil() { {{$l := len .ArgType}}{{$has := false}}	{{if gt $l 0}}{{$ft := index .ArgType 0}}{{if eq $ft "string"}}{{$has = true}}{{end}}{{end}}
		{{if $has}} remote = m.selector.Select(m.svr, "{{$.Service}}", arg0){{else}}remote = m.selector.Select(m.svr, "{{$.Service}}", ""){{end}}
	}
	if remote.IsNil() {
		err = fmt.Errorf("service {{$.Service}} not found")
		return
	}
	call, err := m.svr.AsyncPubWithTimeout(fmt.Sprintf("%s%d:{{$.Name}}.{{.Name}}",m.prefix, remote.ServiceId()), msg.Body, m.timeout)
	msg.Free()
	if err != nil {
		return
	}
	call = <-call.Done
	if call.Err != nil {
		err = call.Err
		return
	}
	{{$l := len .ReturnType}}{{if gt $l 1}}
	for {
		ar := protocol.NewLoadArchive(call.Data)
		{{range $k, $v := .ReturnType}}{{if ne $v "error"}}
		err = ar.Get(&reply{{$k}})
		if err != nil {
			break
		} {{end}} {{end}}
		break
	}{{end}}

	if call.Msg != nil {
		call.Msg.Free()
		call.Msg = nil
	}
	return
}
{{end}}
