package turbogen

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	"github.com/mysll/toolkit"
	"github.com/urfave/cli"
)

var module_tpl = `package {{tolower .Pkg}}
import (
	"context"
	"turboengine/common/utils"
	"turboengine/core/api"
	"turboengine/core/module"
)
 
// Module: 		{{.Name}}
// Auth: 	 	{{.Auth}}
// Data:	  	{{.Time.Format "2006-01-02 15:04:05"}}
// Desc:
type {{.Name}} struct{
	module.Module
}

func (m *{{.Name}}) Name() string{
	return "{{.Name}}"
}

func (m *{{.Name}}) OnPrepare(s api.Service) error {
	m.Module.OnPrepare(s)
	// load module resource
	// load module resource end

	return nil
}

func (m *{{.Name}}) OnStart(ctx context.Context) error {
	m.Module.OnStart(ctx)
	// subscribe subject
	// subscribe subject end
	return nil
}

func (m *{{.Name}}) OnUpdate(t *utils.Time) {

}

func (m *{{.Name}}) OnStop() error {
	return nil
}

`

type ModuleInfo struct {
	Pkg  string
	Name string
	Auth string
	Time time.Time
}

func createModule(pkg, module, auth, path string) {
	t := template.Must(template.New(module).Funcs(template.FuncMap{
		"tolower": strings.ToLower,
	}).Parse(module_tpl))
	outfile := path + "/" + strings.ToLower(module) + ".go"
	if ok, _ := toolkit.PathExists(path); !ok {
		err := os.MkdirAll(path, 0777)
		if err != nil {
			panic(err)
		}
	}

	if ok, _ := toolkit.PathExists(outfile); ok {
		fmt.Print(outfile, " file already exists, overwite it?[y/n]:")
		var confirm string
		fmt.Scanln(&confirm)
		if strings.ToLower(confirm) != "y" {
			fmt.Println("abort")
			return
		}

	}
	f, err := os.Create(outfile)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = t.Execute(f, ModuleInfo{
		Pkg:  pkg,
		Name: module,
		Auth: auth,
		Time: time.Now(),
	})
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("gofmt", "--w", outfile)
	cmd.Run()
	fmt.Println("File created successfully! location: ", outfile)
}

func CreateModule(c *cli.Context) error {
	path := c.String("path")
	if path == "" {
		fmt.Println("miss output path, use --path set output path ")
		return fmt.Errorf("miss output path, use --path set output path ")
	}
	fmt.Print("package name:")
	var pkg string
	fmt.Scanln(&pkg)
	fmt.Print("module name:")
	var sname string
	fmt.Scanln(&sname)
	fmt.Print("auth:")
	var auth string
	fmt.Scanln(&auth)
	createModule(pkg, sname, auth, path)
	return nil
}
