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

	"github.com/vlad-marlo/shortener/internal/config"
	"github.com/vlad-marlo/shortener/internal/grpc"
	_ "github.com/vlad-marlo/shortener/internal/httpserver/middleware"
	"github.com/vlad-marlo/shortener/internal/service"
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
	srvLogger, err := createLogger("service")
	if err != nil {
		panic(fmt.Sprintf("init strorage logger: %v", err))
	}
	defer func() {
		_ = srvLogger.Sync()
	}()

	serverLogger, err := createLogger("server")
	if err != nil {
		panic(fmt.Sprintf("init server logger: %v", err))
	}
	defer func() {
		_ = serverLogger.Sync()
	}()
	if err != nil {
		serverLogger.Fatal(fmt.Sprintf("init config: %v", err))
	}

	storage, err := initStorage(srvLogger)
	if err != nil {
		serverLogger.Fatal(fmt.Sprintf("init storage: %v", err))
	}
	srv := service.New(srvLogger, storage)
	defer func() {
		if err = srv.Close(); err != nil {
			srvLogger.Error("close service", zap.Error(err))
		}
	}()

	// init server
	httpServer := httpserver.New(srv, serverLogger)
	if config.Get().GRPC {
		var grpcServer *grpc.Server
		grpcServer, err = grpc.New(srv, serverLogger)
		if err != nil {
			serverLogger.Fatal("init grpc server", zap.Error(err))
		}
		go func() {
			if err = grpcServer.Start(); err != nil {
				serverLogger.Fatal("grpc server", zap.Error(err))
			}
		}()
		defer func() {
			if err = grpcServer.Close(); err != nil {
				serverLogger.Error("grpc server stop", zap.Error(err))
			}
		}()
	}

	// preparations for graceful shut down
	var sig os.Signal
	interrupt := make(chan os.Signal, 1)
	closed := make(chan struct{})
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		sig = <-interrupt
		if err = httpServer.Close(); err != nil {
			serverLogger.Error(fmt.Sprintf("close server: %v", err))
		}
		close(closed)
	}()

	serverLogger.Info(
		"successfully init server",
		zap.String("bind_addr", config.Get().BindAddr),
		zap.String("storage_type", config.Get().StorageType),
	)

	go func() {
		// logging fatal because listen and server always return not-nil error
		if err = httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
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
func initStorage(logger *zap.Logger) (storage store.Store, err error) {
	cfg := config.Get()
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
