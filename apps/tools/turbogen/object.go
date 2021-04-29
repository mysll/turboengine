package turbogen

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"text/template"

	_ "embed"

	"github.com/mysll/toolkit"
	"github.com/urfave/cli"
)

type AttrDecl struct {
	Name     string
	ArgType  string
	Save     bool
	Public   bool
	Private  bool
	Realtime bool
	Orm      string
}

type ObjectDesc struct {
	Name    string
	Pkg     string
	PkgPath string
	Attrs   []AttrDecl
}

func getTypeAlias(typ string) string {
	if typ == "[]uint8" {
		return "[]byte"
	}
	return typ
}

func getType(typ string) string {
	switch typ {
	case "string":
		return "object.StringHolder"
	case "float32":
		return "object.FloatHolder"
	case "float64":
		return "object.Float64Holder"
	case "int32":
		return "object.IntHolder"
	case "int64":
		return "object.Int64Holder"
	case "Vec2":
		return "object.Vector2Holder"
	case "Vec3":
		return "object.Vector3Holder"
	case "[]uint8":
		return "object.BytesHolder"
	default:
		return "unknown"
	}
}

func getTypeCreate(typ string) string {
	switch typ {
	case "string":
		return "object.NewStringHolder"
	case "float32":
		return "object.NewFloatHolder"
	case "float64":
		return "object.NewFloat64Holder"
	case "int32":
		return "object.NewIntHolder"
	case "int64":
		return "object.NewInt64Holder"
	case "Vec2":
		return "object.NewVector2Holder"
	case "Vec3":
		return "object.NewVector3Holder"
	case "[]uint8":
		return "object.NewBytesHolder"
	default:
		return "object.NewNoneHolder"
	}
}

//go:embed object.tpl
var objectWarp string

func ObjectWrap(s interface{}, pkgpath string, pkg string, path string) {
	ctype := reflect.TypeOf(s)
	count := ctype.Elem().NumField()
	desc := ObjectDesc{}
	desc.PkgPath = pkgpath
	desc.Pkg = pkg
	desc.Name = ctype.Elem().Name()
	desc.Attrs = make([]AttrDecl, 0, count)
	for i := 0; i < count; i++ {
		m := ctype.Elem().Field(i)
		decl := AttrDecl{}
		decl.Name = m.Name
		if m.Type.Kind() != reflect.Slice {
			decl.ArgType = m.Type.Name()
		} else {
			decl.ArgType = "[]" + m.Type.Elem().Name()
		}

		if value, ok := m.Tag.Lookup("attr"); ok {
			values := strings.Split(value, ",")
			for _, v := range values {
				switch v {
				case "save":
					decl.Save = true
				case "public":
					decl.Public = true
				case "private":
					decl.Private = true
				case "realtime":
					decl.Realtime = true
				}
			}
		}
		if value, ok := m.Tag.Lookup("orm"); ok {
			decl.Orm = value
		}

		desc.Attrs = append(desc.Attrs, decl)
	}
	t := template.Must(template.New(desc.Name).Funcs(template.FuncMap{
		"tolower": strings.ToLower,
		"getType": getType,
		"create":  getTypeCreate,
		"alias":   getTypeAlias,
	}).Parse(objectWarp))
	outfile := path + "/" + strings.ToLower(desc.Name) + ".go"
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

//go:embed entity.tpl
var entityDesc string

func CreateEntity(c *cli.Context) error {
	path := c.String("path")
	if path == "" {
		fmt.Println("miss output path, use --path set output path ")
		return fmt.Errorf("miss output path, use --path set output path ")
	}

	fmt.Print("package name:")
	var pkg string
	fmt.Scanln(&pkg)
	fmt.Print("entity name:")
	var name string
	fmt.Scanln(&name)

	makeFile(entityDesc, "def", path, strings.ToLower(name), map[string]interface{}{
		"Name": name,
		"Pkg":  pkg,
	})
	return nil

}
