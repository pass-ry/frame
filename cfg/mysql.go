package cfg

import (
	"sync"

	"github.com/spf13/viper"
	"gitlab.ifchange.com/data/cordwood/mysql"
)

var cfgMySQL mysql.Config

func GetCfgMySQL() mysql.Config {
	cfgMySQLLoaderOnce.Do(new(cfgMySQLLoader).load)
	return cfgMySQL
}

var (
	_                  loader = (*cfgMySQLLoader)(nil)
	cfgMySQLLoaderOnce        = new(sync.Once)
)

type cfgMySQLLoader struct{}

func (l *cfgMySQLLoader) load() {
	viperCfg := viper.GetStringMap("MySQL")
	cfgMySQL = mysql.Config{
		Address:  viperCfg["address"].(string),
		Port:     viperCfg["port"].(string),
		DB:       viperCfg["db"].(string),
		Username: viperCfg["username"].(string),
		Password: viperCfg["password"].(string),
	}
	if keepAlive, ok := viperCfg["keepalive"].(int64); ok {
		cfgMySQL.KeepAlive = int(keepAlive)
	}
	if maxOpenConns, ok := viperCfg["maxopenconns"].(int64); ok {
		cfgMySQL.MaxOpenConns = int(maxOpenConns)
	}
	if maxIdleConns, ok := viperCfg["maxidleconns"].(int64); ok {
		cfgMySQL.MaxIdleConns = int(maxIdleConns)
	}
}
