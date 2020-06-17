package malatd

// SrvConfig
type SrvConfig struct {
	// 监听地址[host:port]
	Address string `yaml:"address"`
	// 是否打开服务运行/请求日志
	RunLog bool `yaml:"run_log"`
}

// NewSrvConfig
func NewSrvConfig() *SrvConfig {
	return &SrvConfig{
		Address: "0.0.0.0:8080",
	}
}
