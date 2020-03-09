package ps

import (
	"context"

	"github.com/google/gops/agent"
	"gitlab.ifchange.com/data/cordwood/log"
)

func Construct(ctx context.Context) {
	go func() {
		<-ctx.Done()
		agent.Close()
		log.Infof("GO-PS is gracefully stopped")
	}()
	log.Infof("GO-PS is starting")
	err := agent.Listen(agent.Options{})
	log.Infof("GO-PS is stopped %v", err)
}
