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

var service_tpl = `package {{tolower .Pkg}} 

import (
	"turboengine/core/api"
	"turboengine/core/service"
)
 
// Service: 	{{.Name}}
// Auth: 	 	{{.Auth}}
// Data:	  	{{.Time.Format "2006-01-02 15:04:05"}}
// Desc:
type {{.Name}} struct{
	service.Service
}

func (s *{{.Name}}) OnPrepare(srv api.Service, args map[string]string) error {
	s.Service.OnPrepare(srv, args)
	// use plugin
	// use plugin end

	// add module
	// add module end 
	
	return nil
}

func (s *{{.Name}}) OnStart() error {
	return nil
}

func (s *{{.Name}}) OnDependReady() {
}

func (s *{{.Name}}) OnShut() bool {
	return true // If you want to close manually return false
}
`

type ServiceInfo struct {
	Pkg  string
	Name string
	Auth string
	Time time.Time
}

func createService(pkg, service, auth, path string) {
	t := template.Must(template.New(service).Funcs(template.FuncMap{
		"tolower": strings.ToLower,
	}).Parse(service_tpl))
	outfile := path + "/" + strings.ToLower(service) + ".go"
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
	err = t.Execute(f, ServiceInfo{
		Pkg:  pkg,
		Name: service,
		Auth: auth,
		Time: time.Now(),
	})
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("gofmt", "--w", outfile)
	cmd.Run()
}

func CreateService(c *cli.Context) error {
	path := c.String("path")
	if path == "" {
		fmt.Println("miss output path, use --path set output path ")
		return fmt.Errorf("miss output path, use --path set output path ")
	}
	fmt.Print("package name:")
	var pkg string
	fmt.Scanln(&pkg)
	fmt.Print("service name:")
	var sname string
	fmt.Scanln(&sname)
	fmt.Print("auth:")
	var auth string
	fmt.Scanln(&auth)

	createService(pkg, sname, auth, path)
	return nil
}
