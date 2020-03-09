package cfg

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
	dfs "gitlab.ifchange.com/data/cordwood/dfs"
)

var cfgDfs dfs.Config

func GetCfgDfs() dfs.Config {
	cfgDfsLoaderOnce.Do(new(cfgDfsLoader).load)
	return cfgDfs
}

var (
	_                loader = (*cfgDfsLoader)(nil)
	cfgDfsLoaderOnce        = new(sync.Once)
)

type cfgDfsLoader struct{}

func (l *cfgDfsLoader) load() {
	viperCfg := viper.GetStringMap("Dfs")
	master, ok := viperCfg["master"].([]interface{})
	if !ok {
		panic(fmt.Errorf("Unknown Dfs-Master %v", viperCfg))
	}
	for _, oneMaster := range master {
		m, ok := oneMaster.(string)
		if ok {
			cfgDfs.Master = append(cfgDfs.Master, m)
		}
	}

	volume, ok := viperCfg["volume"].([]interface{})
	if !ok {
		panic(fmt.Errorf("Unknown Dfs-Volume %v", viperCfg))
	}
	for _, oneVolume := range volume {
		v, ok := oneVolume.(string)
		if ok {
			cfgDfs.Volume = append(cfgDfs.Volume, v)
		}
	}
}
