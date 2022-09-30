package poll

import (
	"log"
	"sync"
	"time"

	"github.com/vlad-marlo/shortener/internal/store"
)

const (
	Interval = 75 * time.Millisecond
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
		ticker: time.NewTicker(Interval),
	}
	return p
}

// DeleteURLs ...
func (p *Poll) DeleteURLs(ids []string, user string) {
	if _, ok := p.input[user]; !ok {
		p.input[user] = make(chan string)
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
			go p.flush()
		case <-p.stop:
			return
		}
	}
}

// flush ...
func (p *Poll) flush() {
	for u, ch := range p.input {
		go func(ch chan string, user string) {
			close(ch)
			var ids []string
			for id := range ch {
				ids = append(ids, id)
			}
			defer close(ch)
			if err := p.store.URLsBulkDelete(ids, user); err != nil {
				log.Printf("urls bulk delete err: %v", err)
			}
		}(ch, u)
	}
}

// Close ...
func (p *Poll) Close() {
	p.stop <- struct{}{}
}
