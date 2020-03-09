package rpc

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"gitlab.ifchange.com/data/cordwood/redis/mock"
)

func Test_RPC(t *testing.T) {
	var (
		pubQueue = "pub_queue"
		job      = []byte("job")

		subQueue = "sub_queue"
		reply    = []byte("reply")
	)

	mockRedis := mock.NewMockRD(gomock.NewController(t))
	mockRedis.EXPECT().ErrNil().Return(errors.New("Err Nil"))
	mockRedis.EXPECT().Do("LPUSH", pubQueue, job).Return(nil, nil)
	mockRedis.EXPECT().DoBytes("RPOP", subQueue).Return(reply, nil)

	err := NewRPC(pubQueue, subQueue).
		OverrideRedis(mockRedis).
		WithCallback(func(r []byte) error {
			if string(r) != string(reply) {
				return errors.Errorf("Want=%v Got=%v", string(reply), string(r))
			}
			return nil
		}).Do(job)

	if err != nil {
		t.Fatal(err)
	}
}
