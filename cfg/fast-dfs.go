package cfg

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
	fastDfs "gitlab.ifchange.com/data/cordwood/fast-dfs"
)

var cfgFastDfs fastDfs.Config

func GetCfgFastDfs() fastDfs.Config {
	cfgFastDfsLoaderOnce.Do(new(cfgFastDfsLoader).load)
	return cfgFastDfs
}

var (
	_                    loader = (*cfgFastDfsLoader)(nil)
	cfgFastDfsLoaderOnce        = new(sync.Once)
)

type cfgFastDfsLoader struct{}

func (l *cfgFastDfsLoader) load() {
	viperCfg := viper.GetStringMap("FastDfs")
	tracker, ok := viperCfg["tracker"].([]interface{})
	if !ok {
		panic(fmt.Errorf("Unknown FastDfs-tracker %v", viperCfg))
	}
	for _, oneTracker := range tracker {
		t, ok := oneTracker.(string)
		if ok {
			cfgFastDfs.Tracker = append(cfgFastDfs.Tracker, t)
		}
	}
	maxConn, ok := viperCfg["maxconn"].(int64)
	if !ok {
		panic(fmt.Errorf("Unknown FastDfs-maxconn %v", viperCfg))
	}
	cfgFastDfs.MaxConn = int(maxConn)
}
