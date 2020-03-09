package rpc

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"gitlab.ifchange.com/data/cordwood/redis"
)

type RPC struct {
	ctx               context.Context
	conn              redis.Client
	pubQueue          string
	subQueue          string
	expectReplyNumber int
	callback          func(reply []byte) error
}

func NewRPC(pubQueue, subQueue string) *RPC {
	return &RPC{
		ctx:               context.Background(),
		conn:              nil,
		pubQueue:          pubQueue,
		subQueue:          subQueue,
		expectReplyNumber: 1,
	}
}

func (rpc *RPC) WithContext(ctx context.Context) *RPC {
	rpc.ctx = ctx
	return rpc
}

func (rpc *RPC) OverrideRedis(conn redis.Client) *RPC {
	rpc.conn = conn
	return rpc
}

func (rpc *RPC) ExpectReplyNum(number int) *RPC {
	rpc.expectReplyNumber = number
	return rpc
}

func (rpc *RPC) WithCallback(callback func(reply []byte) error) *RPC {
	rpc.callback = callback
	return rpc
}

func (rpc *RPC) Do(job []byte) error {
	select {
	case <-rpc.ctx.Done():
		return errors.Errorf("Context Canceled Before Pub")
	default:
	}

	if rpc.conn == nil {
		conn, err := redis.GetConstClient()
		if err != nil {
			return errors.Wrap(err, "redis.GetConstClient")
		}
		rpc.conn = conn
	}

	_, err := rpc.conn.Do("LPUSH", rpc.pubQueue, job)
	if err != nil {
		return errors.Wrapf(err, "LPUSH %s", rpc.pubQueue)
	}

	var sleep time.Duration
	for rpc.expectReplyNumber > 0 {
		select {
		case <-rpc.ctx.Done():
			return errors.Errorf("Context Canceled After Pub")
		case <-time.After(sleep):
			sleep = 0
		}

		reply, err := rpc.conn.DoBytes("RPOP", rpc.subQueue)
		if err == rpc.conn.ErrNil() {
			sleep = time.Millisecond * time.Duration(500)
			continue
		}
		if err != nil {
			return errors.Wrapf(err, "RPOP %s", rpc.subQueue)
		}
		if rpc.callback != nil {
			if err := rpc.callback(reply); err != nil {
				return err
			}
		}
		rpc.expectReplyNumber -= 1
	}
	return nil
}
