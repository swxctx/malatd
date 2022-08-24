package create

const __tpl__ = `// Command __PROJ_NAME__ is the malatd service project.
// The framework reference: https://github.com/swxctx/malatd
package __TPL__

// __API_TPL__ register PULL router
type __API_TPL__ interface {
	V1_Test
}

type V1_Test interface {
	Ping(*PingArgsV1) *PingResultV1
}

type (
	PingArgsV1 = struct{}
	PingResultV1 = struct{}
)
`

const __readme__ = `# ${PROJ_NAME}

${readme}

<br>

*This is a project created by ` + "`malatd gen`" + ` command.*

*[About Malatd Command](https://github.com/swxctx/malatd)*
`

var tplFiles = map[string]string{
	"main.go": `package main

import (
	td "github.com/swxctx/malatd"
	"${import_prefix}/api"
)

func main() {
	// Gen Time: ${project_gen_time}
	srv := td.NewServer(cfg.SrvConfig)
	api.Route(srv, "/${service_api_prefix}")
	srv.Run()
}
`,
	"config.go": `package main

import (
	td "github.com/swxctx/malatd"
	"github.com/swxctx/xlog"
	"github.com/usthooz/gconf"
)

type Config struct {
	SrvConfig *td.SrvConfig ` + "`json:\"srv_config\"`" + `
}

var cfg = &Config{
	SrvConfig: td.NewSrvConfig(),
}

func reload() {
	conf := gconf.NewConf(&gconf.Gconf{
		ConfPath: "./config/config.yaml",
	})

	// get config
	err := conf.GetConf(&cfg)
	if err != nil {
		xlog.Errorf("GetConf Err: %v", err.Error())
	}
}

func init() {
	reload()
}
`,
	"config/config.yaml": "",

	"args/const.gen.go": `package args
${const_list}
`,

	"args/type.gen.go": `package args
import (${import_list}
)
${type_define_list}
`,

	"logic/tmp_code.gen.go": `package logic
import (
	 td "github.com/swxctx/malatd"
	 "${import_prefix}/args"
)
${logic_api_define}
`,

	"api/handler.gen.go": `package api
import (
    td "github.com/swxctx/malatd"
	"github.com/swxctx/malatd/binding"

    "${import_prefix}/logic"
    "${import_prefix}/args"
)
${handler_api_define}
`,

	"api/router.gen.go": `
package api

import (
	td "github.com/swxctx/malatd"
)

func Route(srv *td.Server, rootGroup string) {
	//自定义路由处理
	routeLogic(srv, rootGroup)
	${register_router_list}
}
`,
}
