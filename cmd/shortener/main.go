package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	_ "github.com/vlad-marlo/shortener/internal/httpserver/middleware"
	"github.com/vlad-marlo/shortener/internal/store/filebased"
	"github.com/vlad-marlo/shortener/internal/store/inmemory"
	"github.com/vlad-marlo/shortener/internal/store/sqlstore"

	"github.com/vlad-marlo/shortener/internal/httpserver"
	"github.com/vlad-marlo/shortener/internal/store"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

// main ...
func main() {
	storeLogger, err := createLogger("storage")
	if err != nil {
		panic(fmt.Sprintf("init strorage logger: %v", err))
	}
	defer func() {
		_ = storeLogger.Sync()
	}()

	serverLogger, err := createLogger("server")
	if err != nil {
		panic(fmt.Sprintf("init server logger: %v", err))
	}
	defer func() {
		_ = serverLogger.Sync()
	}()

	config, err := httpserver.NewConfig()
	if err != nil {
		serverLogger.Fatal(fmt.Sprintf("init config: %v", err))
	}

	storage, err := initStorage(config, storeLogger)
	if err != nil {
		serverLogger.Fatal(fmt.Sprintf("init storage: %v", err))
	}
	defer func() {
		if err = storage.Close(); err != nil {
			storeLogger.Fatal(fmt.Sprintf("close storage: %v", err))
		}
	}()

	// init server
	s := httpserver.New(config, storage, serverLogger)

	// preparations for graceful shut down
	var sig os.Signal
	interrupt := make(chan os.Signal, 1)
	closed := make(chan struct{})
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		sig = <-interrupt
		if err = s.Close(); err != nil {
			serverLogger.Error(fmt.Sprintf("close server: %v", err))
		}
		close(closed)
	}()

	serverLogger.Info(
		"successfully init server",
		zap.String("bind_addr", config.BindAddr),
		zap.String("storage_type", config.StorageType),
	)

	go func() {
		// logging fatal because listen and server always return not-nil error
		if err = s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			serverLogger.Fatal(fmt.Sprintf("listen and server server: %v", err))
		}
	}()

	<-closed
	serverLogger.Info(
		"graceful shut down",
		zap.String("signal", sig.String()),
	)
}

// createLogger creates new named logger with stdout and file output.
func createLogger(name string) (*zap.Logger, error) {

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	jsonEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	// textEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	core := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, consoleErrors, highPriority),
		zapcore.NewCore(jsonEncoder, consoleDebugging, lowPriority),
	)

	f := TraceFields()
	f = append(f, zap.String("name", name))
	logger := zap.
		New(core).
		With(f...)
	return logger, nil
}

// initStorage is abstract factory to create new storage object with provided config.
func initStorage(cfg *httpserver.Config, logger *zap.Logger) (storage store.Store, err error) {
	logger.Debug(
		"trace config vars",
		zap.Bool("filename_provided", cfg.FilePath != ""),
		zap.Bool("db_uri_provided", cfg.Database != ""),
		zap.String("storage_type", cfg.StorageType),
	)

	switch cfg.StorageType {
	case store.InMemoryStorage:
		storage = inmemory.New()
	case store.FileBasedStorage:
		storage, err = filebased.New(cfg.FilePath)
	case store.SQLStore:
		storage, err = sqlstore.New(context.Background(), cfg.Database, logger, nil)
	default:
		storage = inmemory.New()
	}
	return
}

// TraceFields return zap fields for logger to trace debug
func TraceFields() []zap.Field {
	return []zap.Field{
		zap.String("build version", buildVersion),
		zap.String("build date", buildDate),
		zap.String("build commit", buildCommit),
	}
}
