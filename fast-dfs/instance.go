package dfs

import (
	"context"
	"runtime"

	"github.com/beiping96/grace"
	"github.com/pkg/errors"
)

var (
	constClient Client
)

func Construct(cfg Config) {
	var err error
	constClient, err = NewClient(&cfg)
	if err != nil {
		panic(errors.Wrapf(err, "New Fast-Dfs Client %+v", cfg))
	}
	grace.Go(clientDestoryMonitor)
}

func clientDestoryMonitor(ctx context.Context) {
	<-ctx.Done()
	runtime.Gosched()
	constClient.Destory()
}

func GetConstClient() Client { return constClient }
