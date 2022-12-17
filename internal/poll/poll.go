package poll

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/vlad-marlo/shortener/internal/store"
)

type (
	Poll struct {
		store  store.Store
		input  chan *task
		logger *zap.Logger
		stop   chan struct{}
	}
	task struct {
		user string
		ids  []string
	}
)

// New ...
func New(store store.Store, logger *zap.Logger) *Poll {
	p := &Poll{
		store:  store,
		input:  make(chan *task, 10),
		stop:   make(chan struct{}),
		logger: logger,
	}
	go p.startPolling()
	return p
}

// DeleteURLs ...
func (p *Poll) DeleteURLs(urls []string, user string) {
	p.logger.Debug(
		"pushing task to queue",
		zap.String("user", user),
		zap.Strings("ids", urls),
	)
	p.input <- &task{
		ids:  urls,
		user: user,
	}
}

func (p *Poll) startPolling() {
	p.logger.Info("starting poller polling")
	for {
		select {
		case <-p.stop:
			return
		case t := <-p.input:
			p.logger.Debug(
				"poll: got new task",
				zap.Strings("ids", t.ids),
				zap.String("user", t.user),
			)
			if err := p.store.URLsBulkDelete(t.ids, t.user); err != nil {
				p.logger.Warn(
					fmt.Sprintf("poll: start_polling: %v", err),
					zap.String("user", t.user),
					zap.Strings("ids", t.ids),
				)
				continue
			}
			p.logger.Debug("successfully done task")
		}
	}
}

func (p *Poll) Close() {
	p.logger.Info("close poller queue")
	close(p.stop)
}
