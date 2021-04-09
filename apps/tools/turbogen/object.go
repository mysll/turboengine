package turbogen

import (
	"fmt"
	"reflect"
	"strings"

	_ "embed"

	"github.com/urfave/cli"
)

type AttrDecl struct {
	Name    string
	ArgType string
}

type ObjectDesc struct {
	Name    string
	Pkg     string
	PkgPath string
	Attrs   []AttrDecl
}

func ObjectWrap(s interface{}, pkgpath string, pkg string, path string) {
	ctype := reflect.TypeOf(s)
	count := ctype.Elem().NumField()
	desc := ObjectDesc{}
	desc.PkgPath = pkgpath
	desc.Pkg = pkg
	desc.Name = ctype.Elem().Name()
	desc.Attrs = make([]AttrDecl, count)
}

//go:embed entity.tpl
var entityDesc string

//go:embed object.tpl
var objectWarp string

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
