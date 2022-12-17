package httpserver

import (
	"sync"

	"go.uber.org/zap"

	"github.com/vlad-marlo/shortener/internal/poll"
	"github.com/vlad-marlo/shortener/internal/store"
	mock_store "github.com/vlad-marlo/shortener/internal/store/mock"
)

var (
	once sync.Once
	c    *Config
	l    *zap.Logger
)

// TestI is giving access to all objects like *testing.B *testing.T in helper functions.
type TestI interface {
	Helper()
}

func TestServer(t TestI, storage store.Store) (*Server, func() error) {
	once.Do(func() {
		c = NewConfig()
		l, _ = zap.NewProduction()
	})
	t.Helper()
	if s, ok := storage.(*mock_store.MockStore); ok {
		s.EXPECT().Close().Return(nil).AnyTimes()
	}
	server := &Server{
		logger: l,
		store:  storage,
		config: c,
		poller: poll.New(storage, l),
	}
	return server, server.Close
}
