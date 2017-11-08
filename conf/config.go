package conf

import (
	"github.com/BurntSushi/toml"
)

var (
	config GlobalConfig
)

type GlobalConfig struct {
	Server serverConfig
	Redis  redisConfig
}

type serverConfig struct {
	TcpBind string `toml:"tcp_bind"`
	RpcBind string `toml:"rpc_bind"`
}

type redisConfig struct {
	Address  string `toml:"address"`
	Password string `toml:"password"`
}

func Parse(cfgPath string) error {
	if _, err := toml.DecodeFile(cfgPath, &config); err != nil {
		return err
	}

	return nil
}

func Config() *GlobalConfig {
	return &config
}
