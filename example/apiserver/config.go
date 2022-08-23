package main

import (
	td "github.com/swxctx/malatd"
	"github.com/swxctx/xlog"
	"github.com/usthooz/gconf"
)

type Config struct {
	SrvConfig *td.SrvConfig `yaml:"srv_config"`
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
