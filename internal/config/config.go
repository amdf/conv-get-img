package config

import (
	"github.com/BurntSushi/toml"
)

//FileName где лежит.
const FileName = "configs/config.toml"

var sc *Config

type ServerCfg struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
}

type GatewayCfg struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
}

type StorageCfg struct {
	Path    string `toml:"path"`
	Timeout string `toml:"timeout"`
}

type Config struct {
	Server  ServerCfg  `toml:"server"`
	Gateway GatewayCfg `toml:"gateway"`
	Storage StorageCfg `toml:"storage"`
}

func Get() *Config {
	return sc
}

//Load from file
func Load() (err error) {
	sc = new(Config)
	_, err = toml.DecodeFile(FileName, sc)

	return err
}

func GetServerAddress() (result string) {
	if nil == sc {
		return
	}
	cfg := sc.Server
	return cfg.Host + ":" + cfg.Port
}

func GetGatewayAddress() (result string) {
	if nil == sc {
		return
	}
	cfg := sc.Gateway
	return cfg.Host + ":" + cfg.Port
}
