package cfg

import (
	"sync"

	"github.com/spf13/viper"
	"gitlab.ifchange.com/data/cordwood/log"
)

var cfgLog log.Config

func GetCfgLog() log.Config {
	cfgLogLoaderOnce.Do(new(cfgLogLoader).load)
	return cfgLog
}

var (
	_                loader = (*cfgLogLoader)(nil)
	cfgLogLoaderOnce        = new(sync.Once)
)

type cfgLogLoader struct{}

func (l *cfgLogLoader) load() {
	viperCfg := viper.GetStringMapString("Log")
	cfgLog = viperCfg["logfile"]
}
