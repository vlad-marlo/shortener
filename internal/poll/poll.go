package poll

import (
	"log"
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
	go p.startPolling()
	return p
}

// DeleteURLs ...
func (p *Poll) DeleteURLs(urls []string, user string) {
	p.input <- &task{
		ids:  urls,
		user: user,
	}
}

func (p *Poll) startPolling() {
	for {
		select {
		case <-p.stop:
			return
		case t := <-p.input:
			go p.deleteURLs(t)
		}
	}
}

func (p *Poll) deleteURLs(t *task) {
	if err := p.store.URLsBulkDelete(t.ids, t.user); err != nil {
		log.Printf("poll: start_polling: %v", err)
	}
}

func (p *Poll) Close() {
	close(p.stop)
}
