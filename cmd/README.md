### malatd
Refer to [tp-micro](https://github.com/xiaoenai/tp-micro/tree/master/cmd/micro)

### Usage

```shell
NAME:
   Malatd project command - A deployment tools of malatd frameware

USAGE:
   malatd [global options] command [command options] [arguments...]

VERSION:
   1.0.0

AUTHOR:
   swxctx

COMMANDS:
   gen      Generate a malatd project
   run      Compile and run go project
   doc      Generate a project README.md(malatd doc || malatd doc -r ${root_group})
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### 模板文件说明
```go
// __API_TPL__ register PULL router
type __API_TPL__ interface {
	V1_Test //接口模块1
	V2_Test //接口模块2
}

type V1_Test interface {//模块接口实现
	Ping(*PingArgsV1) *PingResultV1 //接口
}

type V2_Test interface {//模块接口实现
	PingAdd(*PingArgsV1) *PingResultV1 //接口
}

type (
	PingArgsV1   = struct{
	    A string	
    } //请求参数
	PingResultV1 = struct{} //响应参数
)
```

### 目录结构说明
```
├── README.md
├── __malatd__gen__.lock
├── __malatd__tpl__.go  //模板文件
├── api
│   ├── handler.gen.go  //接口handler
│   ├── handler.go      //用于自定义
│   ├── router.gen.go   //接口router
│   └── router.go       //用于自定义
├── args
│   ├── const.gen.go    //自动常量成成
│   ├── const.go        //自定义常量
│   ├── type.gen.go     //定义结构
│   ├── type.go         //自定义结构
│   └── var.go          //自定义变量
├── config
│   └── config.yaml     //配置文件
├── config.go           //配置加载
├── logic               //接口实现目录
│   └── tmp_code.gen.go
├── main.go             //服务入口文件
└── rerrs               //接口错误定义
└── rerrs.go
```
