package httpserver

import (
	"io"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/vlad-marlo/logger"

	"github.com/vlad-marlo/shortener/internal/poll"
	"github.com/vlad-marlo/shortener/internal/store"
	mock_store "github.com/vlad-marlo/shortener/internal/store/mock"
)

var (
	once sync.Once
	c    *Config
	l    *logrus.Entry
)

func TestServer(t *testing.T, storage store.Store) (*Server, func() error) {
	once.Do(func() {
		c = NewConfig()
		l = logrus.NewEntry(logger.WithOpts(logger.WithOutput(io.Discard)))
	})
	t.Helper()
	if s, ok := storage.(*mock_store.MockStore); ok {
		s.EXPECT().Close().Return(nil).AnyTimes()
	}
	server := &Server{
		logger: l,
		store:  storage,
		config: c,
		poller: poll.New(storage),
	}
	return server, server.Close
}
