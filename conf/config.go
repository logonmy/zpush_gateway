package conf

import (
	"github.com/BurntSushi/toml"
)

var (
	config GlobalConfig
)

type GlobalConfig struct {
	Server   serverConfig   `toml:"server"`
	HTTP     httpConfig     `toml:"http"`
	Redis    redisConfig    `toml:"redis"`
	ZK       zkConfig       `toml:"zookeeper"`
	Influxdb influxdbConfig `toml:"influxdb"`
}

type serverConfig struct {
	NodeId  string `toml:"node"`
	TcpBind string `toml:"tcp_bind"`
	RpcBind string `toml:"rpc_bind"`
}

type httpConfig struct {
	Bind string `toml:"bind"`
}

type redisConfig struct {
	Address  string `toml:"address"`
	Password string `toml:"password"`
}

type zkConfig struct {
	Address string `toml:"address"`
}

type influxdbConfig struct {
	Address  string `toml:"address"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	DB       string `toml:"database"`
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
