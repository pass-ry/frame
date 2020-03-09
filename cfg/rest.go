package cfg

import (
	"sync"

	"github.com/spf13/viper"
)

var cfgREST CfgREST

type CfgREST struct {
	Port string
}

func GetCfgREST() CfgREST {
	cfgRESTLoaderOnce.Do(new(cfgRESTLoader).load)
	return cfgREST
}

var (
	_                 loader = (*cfgRESTLoader)(nil)
	cfgRESTLoaderOnce        = new(sync.Once)
)

type cfgRESTLoader struct{}

func (l *cfgRESTLoader) load() {
	viperCfg := viper.GetStringMapString("REST")
	cfgREST = CfgREST{
		Port: viperCfg["port"],
	}
}
