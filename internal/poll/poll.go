package poll

import (
	"sync"
	"time"

	"github.com/vlad-marlo/shortener/internal/store"
)

const (
	PollInterval = 75 * time.Millisecond
)

type (
	Poll struct {
		store  store.Store
		mu     sync.Mutex
		stop   chan struct{}
		input  map[string]chan string
		ticker *time.Ticker
	}
)

// New ...
func New(store store.Store) *Poll {
	p := &Poll{
		store:  store,
		stop:   make(chan struct{}),
		ticker: time.NewTicker(PollInterval),
	}
	return p
}

// DeleteURLs ...
func (p *Poll) DeleteURLs(ids []string, user string) {
	if ch, _ := p.input[user]; ch == nil {
		p.mu.Lock()
		p.input[user] = make(chan string)
		p.mu.Unlock()
	}

	go func() {
		ch := p.input[user]
		for _, id := range ids {
			ch <- id
		}
	}()
}

func (p *Poll) poll() {
	for {
		select {
		case <-p.ticker.C:
			p.flush()
		case <-p.stop:
			return
		}
	}
}

// flush ...
func (p *Poll) flush() {
	for u, ch := range p.input {
		go func(ch chan string, user string) {
			ids := []string{}
			for id := range ch {
				ids = append(ids, id)
			}
			defer close(ch)
			p.store.URLsBulkDelete(ids, user)
		}(ch, u)
	}
}

// Close ...
func (p *Poll) Close() {
	p.stop <- struct{}{}
}
