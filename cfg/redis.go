package cfg

import (
	"sync"
	"time"

	"github.com/spf13/viper"
	"gitlab.ifchange.com/data/cordwood/redis"
)

var cfgRedis redis.Config

func GetCfgRedis() redis.Config {
	cfgRedisLoaderOnce.Do(new(cfgRedisLoader).load)
	return cfgRedis
}

var (
	_                  loader = (*cfgRedisLoader)(nil)
	cfgRedisLoaderOnce        = new(sync.Once)
)

type cfgRedisLoader struct{}

func (l *cfgRedisLoader) load() {
	viperCfg := viper.GetStringMap("Redis")
	cfgRedis = redis.Config{
		IsCluster: viperCfg["iscluster"].(bool),
		Address:   viperCfg["address"].(string),
	}

	if password, ok := viperCfg["password"].(string); ok {
		cfgRedis.Password = password
	}
	if dialTimeout, ok := viperCfg["dialtimeout"].(int64); ok {
		cfgRedis.DialTimeout = time.Duration(dialTimeout) * time.Second
	}
	if readTimeout, ok := viperCfg["readtimeout"].(int64); ok {
		cfgRedis.ReadTimeout = time.Duration(readTimeout) * time.Second
	}
	if writeTimeout, ok := viperCfg["writetimeout"].(int64); ok {
		cfgRedis.WriteTimeout = time.Duration(writeTimeout) * time.Second
	}
	if prefix, ok := viperCfg["prefix"].(string); ok {
		cfgRedis.Prefix = prefix
	}
}
