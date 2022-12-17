package poll_test

import (
	"go.uber.org/zap"

	"github.com/vlad-marlo/shortener/internal/poll"
)

func Example() {
	// ...
	// init poller
	logger, _ := zap.NewProduction()
	poller := poll.New(nil, logger)
	// always defer poller close
	defer poller.Close()
}
