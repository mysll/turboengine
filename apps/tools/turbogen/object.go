package turbogen

import (
	"fmt"
	"strings"

	_ "embed"

	"github.com/urfave/cli"
)

func ObjectWrap(s interface{}, pkgpath string, pkg string, path string) {
	/*ctype := reflect.TypeOf(s)
	count := ctype.Elem().NumField()*/
}

//go:embed entity.tpl
var entity_desc string

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

	makeFile(entity_desc, "def", path, strings.ToLower(name), map[string]interface{}{
		"Name": name,
		"Pkg":  pkg,
	})
	return nil

}
