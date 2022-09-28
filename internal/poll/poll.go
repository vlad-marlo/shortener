package poll

import (
	"sync"

	"github.com/vlad-marlo/shortener/internal/store"
)

type (
	Poll struct {
		store store.Store
		mu    sync.Mutex
		input chan *task
		stop  chan struct{}
	}
	task struct {
		user string
		ids  []string
	}
)

// New ...
func New(store store.Store) *Poll {
	p := &Poll{
		store: store,
		input: make(chan *task),
		stop:  make(chan struct{}),
	}
	return p
}

func (p *Poll) DeleteURLs(urls []string, user string) {
	p.input <- &task{
		ids:  urls,
		user: user,
	}
}
