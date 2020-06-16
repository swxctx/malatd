package malatd

// SrvConfig
type SrvConfig struct {
	// 监听地址[host:port]
	Address string `yaml:"address"`
}

// NewSrvConfig
func NewSrvConfig() *SrvConfig {
	return &SrvConfig{
		Address: "0.0.0.0:8080",
	}
}
