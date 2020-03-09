package cfg

import (
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var cfgCustom CfgCustom

type CfgCustom struct {
	storage map[string]string
}

func GetCfgCustom() CfgCustom {
	cfgCustomLoaderOnce.Do(new(cfgCustomLoader).load)
	return cfgCustom
}

func (cfg CfgCustom) GetAll() map[string]string {
	if cfg.storage == nil {
		panic("Custom Config Storage is not exist")
	}
	return cfg.storage
}

func (cfg CfgCustom) Get(key string) (value string) {
	if cfg.storage == nil {
		panic("Custom Config Storage is not exist")
	}
	value = cfg.storage[strings.ToLower(key)]
	return
}

func (cfg CfgCustom) GetWithOK(key string) (value string, ok bool) {
	if cfg.storage == nil {
		panic("Custom Config Storage is not exist")
	}
	value, ok = cfg.storage[strings.ToLower(key)]
	return
}

var (
	_                   loader = (*cfgCustomLoader)(nil)
	cfgCustomLoaderOnce        = new(sync.Once)
)

type cfgCustomLoader struct{}

func (l *cfgCustomLoader) load() {
	cfgCustom.storage = viper.GetStringMapString("Custom")
}
