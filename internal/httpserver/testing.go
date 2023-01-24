package httpserver

import (
	"fmt"
	"sync"

	"go.uber.org/zap"

	"github.com/vlad-marlo/shortener/internal/config"
	srv "github.com/vlad-marlo/shortener/internal/service"
	"github.com/vlad-marlo/shortener/internal/store"
	mock_store "github.com/vlad-marlo/shortener/internal/store/mock"
)

var (
	// once ...
	once sync.Once
	// o ...
	c *config.Config
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
		cfg := zap.Config{
			Level:            zap.NewAtomicLevelAt(zap.PanicLevel),
			Development:      true,
			Encoding:         "console",
			EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
			OutputPaths:      []string{},
			ErrorOutputPaths: []string{},
		}
		l = zap.Must(cfg.Build())
		c = config.Get()
	})
	t.Helper()
	if s, ok := storage.(*mock_store.MockStore); ok {
		s.EXPECT().Close().Return(nil).AnyTimes()
	}
	server := &Server{
		logger: l,
		srv:    srv.New(l, storage),
		config: c,
	}
	closer, ok := server.srv.(interface {
		Close() error
	})
	if ok {
		return server, func() error {
			// return combined error if any error will be returned in chain
			err1 := closer.Close()
			err2 := server.Close()
			if err1 == nil {
				if err2 == nil {
					return nil
				}
				return err2
			}
			if err2 == nil {
				return err1
			}
			return fmt.Errorf("service error: %v; server error: %w", err1, err2)
		}
	}
	return server, server.Close
}
