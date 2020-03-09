package dfs

import (
	"github.com/pkg/errors"
)

var (
	constConfig Config
)

func Construct(cfg Config) {
	if len(cfg.Master) == 0 || len(cfg.Volume) == 0 {
		panic(errors.Errorf("Construct Dfs Error Nil Config %+v",
			cfg))
	}
	constConfig = cfg
}
