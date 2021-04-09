package turbogen

import (
	_ "embed"
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

//go:embed service.tpl
var service_tpl string

var server_rpc = `package proto
var reg = make(map[string]interface{})
`

//go:embed proto_test.tpl
var server_rpc_test string

var server_entity = `package def
var entities = make(map[string]interface{})
`

//go:embed entity_test.tpl
var server_entity_test string

//go:embed service_main.tpl
var server_main string

//go:embed config.tpl
var config string

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
	tpl := fmt.Sprintf(config, name)

	if err := ioutil.WriteFile(fmt.Sprintf("%s/%s.toml", path, strings.ToLower(name)), []byte(tpl), 0777); err != nil {
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

	entity_def := path + "/internal/def"
	if ok, _ := toolkit.PathExists(entity_def); !ok {
		err := os.MkdirAll(entity_def, 0777)
		if err != nil {
			panic(err)
		}
	}

	entity := path + "/internal/entity"
	if ok, _ := toolkit.PathExists(entity); !ok {
		err := os.MkdirAll(entity, 0777)
		if err != nil {
			panic(err)
		}
	}

	info := ServiceInfo{
		Pkg:  pkg,
		Name: name,
		Auth: auth,
		Time: time.Now(),
	}
	makeFile(service_tpl, name, fmt.Sprintf("%s/%s", path, pkg), strings.ToLower(name), info)

	makeFile(server_main, "main", path, "main", info)

	makeFile(server_rpc, "proto", proto, "proto", info)
	makeFile(server_rpc_test, "proto_test", proto, "proto_test", info)

	makeFile(server_entity, "entity", entity_def, "entity", info)
	makeFile(server_entity_test, "entity_test", entity_def, "entity_test", info)
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
