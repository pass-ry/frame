package cfg

import (
	"sync"

	"github.com/spf13/viper"
)

var cfgRPC CfgRPC

type CfgRPC struct {
	Port string
}

func GetCfgRPC() CfgRPC {
	cfgRPCLoaderOnce.Do(new(cfgRPCLoader).load)
	return cfgRPC
}

var (
	_                loader = (*cfgRPCLoader)(nil)
	cfgRPCLoaderOnce        = new(sync.Once)
)

type cfgRPCLoader struct{}

func (l *cfgRPCLoader) load() {
	viperCfg := viper.GetStringMapString("RPC")
	cfgRPC = CfgRPC{
		Port: viperCfg["port"],
	}
}
