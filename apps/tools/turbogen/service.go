package turbogen

import (
	"fmt"
	"io/ioutil"
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
	coreapi "turboengine/core/api"
	"turboengine/core/service"
)
 
// Service: 	{{.Name}}
// Auth: 	 	{{.Auth}}
// Data:	  	{{.Time.Format "2006-01-02 15:04:05"}}
// Desc:
type {{.Name}} struct{
	service.Service
}

func (s *{{.Name}}) OnPrepare(srv coreapi.Service, args map[string]string) error {
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

var server_rpc = `package proto
var reg = make(map[string]interface{})
`

var server_rpc_test = `package proto
import(
	"testing"
	"turboengine/apps/tools/turbogen"
)

func TestCreate(t *testing.T) {
	for _, v := range reg {
		turbogen.Generate(v, "rpc", "../rpc")
	}
}

`

var server_main = `package main
import (
	"turboengine/common/log"
	"turboengine/core/service"
	"./{{.Pkg}}"
)

func main() {
	log.Init(nil)
	defer log.Close()

	cfg := new(service.Config)
	if err := cfg.LoadFromToml("./conf/main.toml"); err != nil {
		panic(err)
	}
	{{tolower .Name}} := service.New(new({{.Pkg}}.{{.Name}}), cfg)
	{{tolower .Name}}.Start()
	{{tolower .Name}}.Wait()
}
`

type ServiceInfo struct {
	Pkg  string
	Name string
	Auth string
	Time time.Time
}

func makeFile(tpl string, name, path, file string, data interface{}) {
	t := template.Must(template.New(name).Funcs(template.FuncMap{
		"tolower": strings.ToLower,
	}).Parse(tpl))

	outfile := fmt.Sprintf("%s/%s.go", path, file)
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
	err = t.Execute(f, data)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("gofmt", "--w", outfile)
	cmd.Run()
	fmt.Println("File created successfully! location: ", outfile)
}

func saveConfig(name string, path string) {
	tpl := fmt.Sprintf(`
# service id
ID = 1
# service name
Name = "%s"
# nats agent url
NatsUrl = "nats://0.0.0.0:4222"
# if expose is true, it's a outgoing service.
Expose = false
Addr = "0.0.0.0"
Port = 0
FPS = 30

# depend other service
# [[Depend]]
# Name = ""
# Count = 1

# custom args
[Args]
# Key = "value"
`, name)

	if err := ioutil.WriteFile(path+"/main.toml", []byte(tpl), 0777); err != nil {
		panic(err)
	}
}

func createService(pkg, name, auth, path string) {

	conf := path + "/conf"
	if ok, _ := toolkit.PathExists(conf); !ok {
		err := os.MkdirAll(conf, 0777)
		if err != nil {
			panic(err)
		}
	}
	saveConfig(name, conf)

	rpc := path + "/api/rpc"
	if ok, _ := toolkit.PathExists(rpc); !ok {
		err := os.MkdirAll(rpc, 0777)
		if err != nil {
			panic(err)
		}
	}

	proto := path + "/api/proto"
	if ok, _ := toolkit.PathExists(proto); !ok {
		err := os.MkdirAll(proto, 0777)
		if err != nil {
			panic(err)
		}
	}

	mod := path + "/mod"
	if ok, _ := toolkit.PathExists(mod); !ok {
		err := os.MkdirAll(mod, 0777)
		if err != nil {
			panic(err)
		}
	}

	makeFile(service_tpl, name, fmt.Sprintf("%s/%s", path, pkg), strings.ToLower(name),
		ServiceInfo{
			Pkg:  pkg,
			Name: name,
			Auth: auth,
			Time: time.Now(),
		})

	makeFile(server_main, "main", path, "main", ServiceInfo{
		Pkg:  pkg,
		Name: name,
		Auth: auth,
		Time: time.Now(),
	})

	makeFile(server_rpc, "proto", proto, "proto", nil)
	makeFile(server_rpc_test, "proto_test", proto, "proto_test", nil)
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
