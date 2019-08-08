package turbogen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/mysll/toolkit"
	"github.com/urfave/cli"
)

var tpl = `
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
	ar := protocol.NewLoadArchiver(data)
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
	if err := svr.Sub(fmt.Sprintf("%s%d.{{$.Name}}.{{.Name}}", prefix, svr.ID()), m.{{.Name}}); err != nil {
		return err
	}{{end}}
	return nil
}

// client
type {{$.Name}}_RPC_Go_{{$.Ver}}_Client struct {
	svr coreapi.Service
	prefix string
	timeout time.Duration
}

func (m *{{$.Name}}_RPC_Go_{{$.Ver}}_Client) Redirect(svr coreapi.Service) {
	m.svr = svr
}

{{range .Methods}}
// {{.Name}} must call in a new goroutine, if call in service's goroutine, it will be dead lock
func (m *{{$.Name}}_RPC_Go_{{$.Ver}}_Client) {{.Name}}({{range $k, $v := .ArgType}}{{if ne $k 0}},{{end}}arg{{$k}} {{$v}}{{end}}) ({{range $k, $v := .ReturnType}}{{if ne $k 0}},{{end}}{{if eq $v "error"}}err{{else}}reply{{$k}}{{end}} {{$v}}{{end}}) {
	sr := protocol.NewAutoExtendArchive(128)
	{{range $k, $v := .ArgType}}err = sr.Put(arg{{$k}})
	if err != nil {
		return
	}
	{{end}}
	msg := sr.Message()
	call, err := m.svr.PubWithTimeout(fmt.Sprintf("%s%d.{{$.Name}}.{{.Name}}",m.prefix, m.svr.ID()), msg.Body, m.timeout)
	msg.Free()
	if err != nil {
		return
	}
	call.Done = make(chan *coreapi.Call, 1)
	call = <-call.Done
	if call.Err != nil {
		err = call.Err
		return
	}
	{{$l := len .ReturnType}}{{if gt $l 1}}
	for {
		ar := protocol.NewLoadArchiver(call.Data)
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

func New{{.Name}}Consumer(svr coreapi.Service, prefix string, timeout time.Duration) *proto.{{.Name}} {
	m := new(proto.{{.Name}})
	mc := new({{$.Name}}_RPC_Go_{{$.Ver}}_Client)
	mc.svr = svr
	mc.prefix = prefix
	mc.timeout = timeout
	m.XXX = mc	{{range .Methods}}
	m.{{.Name}}=mc.{{.Name}}{{end}}
	return m
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
	timeout time.Duration
	handler I{{$.Name}}_RPC_Go_{{$.Ver}}_Handler
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
	call, err := m.svr.PubWithTimeout(fmt.Sprintf("%s%d.{{$.Name}}.{{.Name}}",m.prefix, m.svr.ID()), msg.Body, m.timeout)
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
		ar := protocol.NewLoadArchiver(call.Data)
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

func New{{.Name}}ConsumerWithHandle(svr coreapi.Service, prefix string, timeout time.Duration, handler I{{$.Name}}_RPC_Go_{{$.Ver}}_Handler) *proto.{{.Name}} {
	m := new(proto.{{.Name}})
	mc := new({{$.Name}}_RPC_Go_{{$.Ver}}_Client_Handle)
	mc.svr = svr
	mc.prefix = prefix
	mc.timeout = timeout
	mc.handler = handler
	m.XXX = mc	{{range .Methods}}
	m.{{.Name}}=mc.{{.Name}}{{end}}
	return m
}
`

type FuncDecl struct {
	Name       string
	ArgType    []string
	ReturnType []string
}

type RpcDesc struct {
	Name    string
	Ver     string
	Pkg     string
	PkgPath string
	Imp     map[string]struct{}
	Methods []FuncDecl
}

func TypeName(t reflect.Type) (pkg string, typ string) {
	ptr := ""
	if t.Kind() == reflect.Ptr {
		ptr = "*"
		t = t.Elem()
	}

	pkg = t.PkgPath()
	var base string
	if pkg != "" {
		base = filepath.Base(pkg) + "."
	}

	typ = ptr + base + t.Name()
	return
}

func Generate(s interface{}, pkgpath string, pkg string, path string) {
	ctype := reflect.TypeOf(s)
	count := ctype.Elem().NumField()
	desc := RpcDesc{}
	desc.Ver = "V1_0"
	desc.PkgPath = pkgpath
	desc.Pkg = pkg
	desc.Name = ctype.Elem().Name()
	desc.Imp = make(map[string]struct{})
	desc.Methods = make([]FuncDecl, 0, count)
	for i := 0; i < count; i++ {
		m := ctype.Elem().Field(i)
		if m.Name == "Ver" {
			v := m.Tag.Get("version")
			desc.Ver = "V" + strings.Replace(v, ".", "_", -1)
			continue
		}

		mtype := m.Type
		if mtype.Kind() != reflect.Func {
			continue
		}

		decl := FuncDecl{}
		decl.Name = m.Name
		decl.ArgType = make([]string, mtype.NumIn())
		decl.ReturnType = make([]string, mtype.NumOut())
		for i := 0; i < mtype.NumIn(); i++ {
			pkg, typ := TypeName(mtype.In(i))
			if pkg != "" {
				desc.Imp[pkg] = struct{}{}
			}
			decl.ArgType[i] = typ
		}
		for i := 0; i < mtype.NumOut(); i++ {
			pkg, typ := TypeName(mtype.Out(i))
			if pkg != "" {
				desc.Imp[pkg] = struct{}{}
			}
			decl.ReturnType[i] = typ
		}

		desc.Methods = append(desc.Methods, decl)
	}

	t := template.Must(template.New(desc.Name).Parse(tpl))
	outfile := path + "/" + strings.ToLower(desc.Name) + "_rpc_wrap.go"
	if ok, _ := toolkit.PathExists(path); !ok {
		err := os.MkdirAll(path, 0777)
		if err != nil {
			panic(err)
		}
	}

	f, err := os.Create(outfile)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = t.Execute(f, desc)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("gofmt", "--w", outfile)
	cmd.Run()
}

var tpl_proto = `package {{.Pkg}}

type {{.Name}} struct {
	Ver   string {{.Tag}}
	XXX   interface{}
	// custom method begin

	// custom method end
}

func init() {
	reg["{{.Name}}"] = new({{.Name}})
}
`

func CreateProto(c *cli.Context) error {
	path := c.String("path")
	if path == "" {
		fmt.Println("miss output path, use --path set output path ")
		return fmt.Errorf("miss output path, use --path set output path ")
	}

	fmt.Print("package name:")
	var pkg string
	fmt.Scanln(&pkg)
	fmt.Print("proto name:")
	var name string
	fmt.Scanln(&name)
	fmt.Print("auth:")
	var auth string
	fmt.Scanln(&auth)

	makeFile(tpl_proto, "proto", path, strings.ToLower(name), map[string]interface{}{
		"Name": name,
		"Pkg":  pkg,
		"Tag":  "`version:\"1.0.0\"`",
	})
	return nil

}
