package httpserver

import (
	"sync"

	"go.uber.org/zap"

	"github.com/vlad-marlo/shortener/internal/poll"
	"github.com/vlad-marlo/shortener/internal/store"
	mock_store "github.com/vlad-marlo/shortener/internal/store/mock"
)

var (
	// once ...
	once sync.Once
	// o ...
	c *Config
	// l ...
	l *zap.Logger
)

// TestI is giving access to all objects like *testing.B *testing.T in helper functions.
type TestI interface {
	Helper()
	Fatalf(format string, args ...any)
}

// TestServer returns server instance, prepared for testing. Always defer func which is returned by TestServer.
func TestServer(t TestI, storage store.Store) (*Server, func() error) {
	once.Do(func() {
		var err error
		c, err = NewConfig()
		if err != nil {
			t.Fatalf("init test config: %v", err)
		}
		cfg := zap.Config{
			Level:            zap.NewAtomicLevelAt(zap.PanicLevel),
			Development:      true,
			Encoding:         "console",
			EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
			OutputPaths:      []string{},
			ErrorOutputPaths: []string{},
		}
		l = zap.Must(cfg.Build())
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
