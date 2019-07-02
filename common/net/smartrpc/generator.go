package smartrpc

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
)

var tpl = `
package {{.Pkg}}

import(
	"turboengine/common/net/rpc"
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
type {{$.Name}}_RPC_Go_{{$.Ver}}_{{.Name}} struct {
	{{range $k, $v := .ArgType}}
	Arg{{$k}} {{$v}}{{end}}
}

type {{$.Name}}_RPC_Go_{{$.Ver}}_{{.Name}}_Reply struct {
	{{range $k, $v := .ReturnType}}{{if ne $v "error"}}
	Arg{{$k}} {{$v}}{{end}}{{end}}
}
{{end}}
{{range .Methods}}
func (math *{{$.Name}}_RPC_Go_{{$.Ver}}) {{.Name}}(arg *{{$.Name}}_RPC_Go_{{$.Ver}}_{{.Name}}, reply *{{$.Name}}_RPC_Go_{{$.Ver}}_{{.Name}}_Reply) (err error) {
	{{range $k, $v := .ReturnType}}{{if ne $v "error"}}reply.Arg{{$k}},{{end}}{{end}}err = math.handler.{{.Name}}({{range $k, $v := .ArgType}}{{if ne $k 0}},{{end}}arg.Arg{{$k}}{{end}})
	return
}
{{end}}
func Set{{.Name}}Provider(svr *rpc.Server, name string, provider I{{.Name}}_RPC_Go_{{.Ver}}) {
	m := new({{.Name}}_RPC_Go_{{.Ver}})
	m.handler = provider
	regName := "{{.Name}}"
	if name != "" {
		regName = name
	}
	svr.RegisterName(regName + "_{{.Ver}}", m)
}

// client
type {{$.Name}}_RPC_Go_{{$.Ver}}_Client struct {
	c   *rpc.Client
	srv string
}

func (m *{{$.Name}}_RPC_Go_{{$.Ver}}_Client) Redirect(c *rpc.Client) {
	m.c = c
}

{{range .Methods}}
func (m *{{$.Name}}_RPC_Go_{{$.Ver}}_Client) {{.Name}}({{range $k, $v := .ArgType}}{{if ne $k 0}},{{end}}arg{{$k}} {{$v}}{{end}}) ({{range $k, $v := .ReturnType}}{{if ne $k 0}},{{end}}{{$v}}{{end}}) {
	_arg := &{{$.Name}}_RPC_Go_{{$.Ver}}_{{.Name}}{}
	{{range $k, $v := .ArgType}}_arg.Arg{{$k}}=arg{{$k}}
	{{end}}
	_reply := &{{$.Name}}_RPC_Go_{{$.Ver}}_{{.Name}}_Reply{}
	err := m.c.Call(m.srv+"_V1_0.{{.Name}}", _arg, _reply)
	{{$l := len .ReturnType}}{{if gt $l 1}}return _reply.Arg0, err{{else}}	return err{{end}}
}
{{end}}
func New{{.Name}}Consumer(client *rpc.Client, srv string) *{{.Name}} {
	m := new({{.Name}})
	mc := new({{$.Name}}_RPC_Go_{{$.Ver}}_Client)
	mc.c = client
	mc.srv = "{{.Name}}"
	if srv != "" {
		mc.srv = srv
	}
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

func Generate(s interface{}, pkg string, path string) {
	ctype := reflect.TypeOf(s)
	count := ctype.Elem().NumField()
	desc := RpcDesc{}
	desc.Ver = "V1_0"
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
