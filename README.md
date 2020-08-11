# turboengine
turbo engine

# quick start
示例
## 创建一个Echo服务:
##    
    cd apps/tools
    go run main.go create service --path ../echo
    输入包名: echo
    输入服务名： Echo  （注意服务器用大写开头）
完成后，会在上一级目录创建Echo服务的基础结构和代码：
## 
    echo
    ├─api           
    │  ├─proto      rpc协议
    │  └─rpc        rpc包装类存放的地方，由工具自动生成
    ├─conf          配置文件目录
    ├─echo          Echo服务
    ├─internal      内部使用的共用代码放在这里
    └─mod           服务需要的module放在这里
## 创建rpc
proto目录下新建echo.go,内容如下：
##  
    type Echo struct {
        Ver string `version:"1.0.0"`
        XXX interface{}
        // custom method begin
        Echo func(string) (string, error)
        // custom method end
    }

    func init() {
        reg["Echo"] = new(Echo)
    }
定义了一个Echo方法，接收一个字符串，原样返回字符串.  
进入proto目录,运行proto_test.go里面的方法TestCreate，将在rpc目录下，生成包装类echo_rpc_wrap.go

## 创建module
##  
    go run main.go create module --path ../echo/mod/echo
    package name:echo
    module name:Echo
执行后将在echo/mod/echo下新建一个echo.go  
实现rpc的Echo方法, 在echo/mod/echo目录下新建echo_rpc.go：
##  
    package echo

    type EchoServer struct {
    }

    func (e *EchoServer) Login(input string) (string, error) {
        return input, nil
    }
回到echo/mod/echo/echo.go  
改写如下方法：
##  
    func (m *Echo) OnStart(ctx context.Context) error {
        m.Module.OnStart(ctx)
        // subscribe subject
        rpc.SetEchoProvider(m.Srv, "", &EchoServer{}) // 关联rpc接口与实现
        // subscribe subject end
        return nil
    }

## service关联module
打开echo/echo/echo.go
##  
    func (s *Echo) OnPrepare(srv coreapi.Service, args map[string]string) error {
        s.Service.OnPrepare(srv, args)
        // use plugin
        // use plugin end

        // add module
        m := module.New(&echo.Echo{}, false)
	    srv.AddModule(m)
        // add module end

        return nil
    }