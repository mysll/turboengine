# turboengine
turbo engine

# quick start
示例创建一个Echo服务:
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
创建rpc,proto目录下新建echo.go,内容如下：
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
定义了一个Echo方法，接收一个字符串，原样返回字符串
进入proto目录,运行proto_test.go里面的方法TestCreate，将在rpc目录下，生成包装类echo_rpc_wrap.go


