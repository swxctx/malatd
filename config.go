package malatd

// SrvConfig
type SrvConfig struct {
	// 监听地址[host:port]
	ListenAddress string `yaml:"listen_address"`
}

// NewSrvConfig
func NewSrvConfig()*SrvConfig {
	return &SrvConfig{
		ListenAddress: "0.0.0.0:8080",
	}
}