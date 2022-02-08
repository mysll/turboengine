package turbogen

import (
	_ "embed"
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

//go:embed proto.tpl
var tpl string

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
	Service string
}

func TypeName(t reflect.Type) (pkg string, typ string) {
	ptr := ""
	if t.Kind() == reflect.Ptr {
		ptr = "*"
		t = t.Elem()
	}
	if t.Kind() == reflect.Slice {
		ptr = "[]"
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

func Generate(s any, pkgpath string, pkg string, path string, service string) {
	ctype := reflect.TypeOf(s)
	count := ctype.Elem().NumField()
	desc := RpcDesc{}
	desc.Ver = "V1_0"
	desc.PkgPath = pkgpath
	desc.Pkg = pkg
	desc.Name = ctype.Elem().Name()
	desc.Imp = make(map[string]struct{})
	desc.Methods = make([]FuncDecl, 0, count)
	desc.Service = service
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

//go:embed proto_struct.tpl
var tpl_proto string

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

	makeSourceFile(tpl_proto, "proto", path, strings.ToLower(name), map[string]any{
		"Name": name,
		"Pkg":  pkg,
		"Tag":  "`version:\"1.0.0\"`",
	})
	return nil

}
