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

func New{{.Name}}Consumer(svr coreapi.Service, prefix string, dest protocol.Mailbox, selector coreapi.Selector, timeout time.Duration) *proto.{{.Name}} {
	m := new(proto.{{.Name}})
	mc := new({{$.Name}}_RPC_Go_{{$.Ver}}_Client)
	mc.svr = svr
	mc.dest = dest
	mc.prefix = prefix
	mc.timeout = timeout
	mc.selector=selector
	m.XXX = mc	{{range .Methods}}
	m.{{.Name}}=mc.{{.Name}}{{end}}
	return m
}

func New{{.Name}}ConsumerBySelector(svr coreapi.Service, prefix string, selector coreapi.Selector, timeout time.Duration) *proto.{{.Name}} {
	return New{{.Name}}Consumer(svr, prefix, 0, selector, timeout)
}

func New{{.Name}}ConsumerByMailbox(svr coreapi.Service, prefix string, remote protocol.Mailbox, timeout time.Duration) *proto.{{.Name}} {
	return New{{.Name}}Consumer(svr, prefix, remote, nil, timeout)
}

{{range .Methods}}
type {{$.Name}}_RPC_Go_{{$.Ver}}_{{.Name}}_Reply struct {
	{{range $k, $v := .ReturnType}}{{if ne $v "error"}}
	Arg{{$k}} {{$v}}{{end}}{{end}}
}
{{end}}

type I{{$.Name}}_RPC_Go_{{$.Ver}}_Handler interface {
	{{range .Methods}}
	On{{.Name}}({{range $k, $v := .ReturnType}}{{if ne $k 0}},{{end}}{{$v}}{{end}}){{end}}
}

type {{$.Name}}_RPC_Go_{{$.Ver}}_Client_Handle struct {
	svr     coreapi.Service
	prefix  string
	dest    protocol.Mailbox
	timeout time.Duration
	handler I{{$.Name}}_RPC_Go_{{$.Ver}}_Handler
	selector coreapi.Selector
}

func (m *{{$.Name}}_RPC_Go_{{$.Ver}}_Client_Handle) Redirect(dest protocol.Mailbox) {
	m.dest = dest
}

func (m *{{$.Name}}_RPC_Go_{{$.Ver}}_Client_Handle)  SetSelector(selector coreapi.Selector) {
	m.selector = selector
}

{{range .Methods}}
func (m *{{$.Name}}_RPC_Go_{{$.Ver}}_Client_Handle) {{.Name}}({{range $k, $v := .ArgType}}{{if ne $k 0}},{{end}}arg{{$k}} {{$v}}{{end}}) ({{range $k, $v := .ReturnType}}{{if ne $k 0}},{{end}}{{if eq $v "error"}}err{{else}}reply{{$k}}{{end}} {{$v}}{{end}}) {
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
	call, err := m.svr.PubWithTimeout(fmt.Sprintf("%s%d:{{$.Name}}.{{.Name}}",m.prefix, remote.ServiceId()), msg.Body, m.timeout)
	msg.Free()
	if err != nil {
		return
	}
	call.Callback = m.On{{.Name}}
	return
}

func (m *{{$.Name}}_RPC_Go_{{$.Ver}}_Client_Handle) On{{.Name}}(call *coreapi.Call) {
	{{$l := len .ReturnType}}{{if gt $l 1}}var reply {{$.Name}}_RPC_Go_{{$.Ver}}_{{.Name}}_Reply{{end}}
	var err error
	err = call.Err
	if err != nil {
		m.handler.On{{.Name}}({{range $k, $v := .ReturnType}}{{if ne $v "error"}}reply.Arg{{$k}},{{end}}{{end}}err)
		return
	}
	{{if gt $l 1}}
	for {
		ar := protocol.NewLoadArchive(call.Data)
		{{range $k, $v := .ReturnType}}{{if ne $v "error"}}
		err = ar.Get(&reply.Arg{{$k}})
		if err != nil {
			break
		} {{end}} {{end}}
		break
	}{{end}}
	m.handler.On{{.Name}}({{range $k, $v := .ReturnType}}{{if ne $v "error"}}reply.Arg{{$k}},{{end}}{{end}}err)
}

{{end}}
