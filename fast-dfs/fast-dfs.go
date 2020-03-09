package dfs

import (
	fastDfsCli "github.com/beiping96/fdfs"
	"github.com/pkg/errors"
)

type Config = fastDfsCli.Config

type Client interface {
	Get(remoteFileID string) (data []byte, err error)
	Set(ext string, data []byte) (remoteFileID string, err error)
	Del(remoteFileID string) (err error)

	Destory()
}

func NewClient(cfg *Config) (cli Client, err error) {
	sourceCli, err := fastDfsCli.NewClient(cfg)
	if err != nil {
		return nil, errors.Wrapf(err,
			"github.com/beiping96/fdfs.NewClient(%+v)", cfg)
	}
	return &client{sourceCli}, nil
}

var _ Client = (*client)(nil)

type client struct{ *fastDfsCli.Client }

func (cli *client) Get(remoteFileID string) (data []byte, err error) {
	return cli.Client.DownloadToBuffer(remoteFileID, 0, 0)
}

func (cli *client) Set(ext string, data []byte) (remoteFileID string, err error) {
	return cli.Client.UploadByBuffer(data, ext)
}

func (cli *client) Del(remoteFileID string) (err error) {
	return cli.Client.DeleteFile(remoteFileID)
}
