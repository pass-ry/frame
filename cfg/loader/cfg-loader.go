/*
WARNING!!!

This package supports this following to load config:
0. load config in dev

1. load config by set ENV (in OS);
2. load config by declare ENV;
3. load config by declare config file.

The Sum of all ways are executed is ONE.

There will PANIC after got error.

Up to now, Hot-Reload is unsupported.
*/
package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

var (
	once     = new(sync.Once)
	ENV  Env = "nil"
)

func LoadCfgByEnv() {
	once.Do(loadCfgByEnv)
}

func LoadCfgByDeclareEnv(env Env) {
	once.Do(func() { loadCfgByDeclareEnv(env) })
}

func LoadCfgByPath(path string) {
	once.Do(func() { loadCfgByPath(path) })
}

type Env = string

const (
	DEV  Env = "dev"
	TEST Env = "test"
	PROD Env = "prod"
)

func loadCfgByEnv() {
	err := viper.BindEnv("ENV")
	if err != nil {
		panic(fmt.Errorf("Bind ENV %v",
			err))
	}
	var env Env = viper.GetString("ENV")

	loadCfgByDeclareEnv(env)
}

func loadCfgByDeclareEnv(env Env) {
	if len(env) == 0 {
		panic(fmt.Errorf(`Not found ENV:
		 Not support ENV=%s
		 Declare ENV as dev, test or prod.`,
			env))
	}
	ENV = env

	path := filepath.Join(".", "config", env)
	loadCfgByPath(path)
}

func loadCfgByPath(path string) {
	viper.SetConfigName("config")
	viper.AddConfigPath(path)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("load config error %v",
			err))
	}
}

func LoadCfgInDev(project string) {
	goPath := os.Getenv("GOPATH")
	if len(goPath) == 0 {
		panic("GOPATH not exported")
	}
	cfgPath := filepath.Join(goPath, "src", project, "config", "dev")
	loadCfgByPath(cfgPath)
}
