package main

import (
	"context"
	"io"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/vlad-marlo/logger"
	"github.com/vlad-marlo/logger/hook"

	_ "github.com/vlad-marlo/shortener/internal/httpserver/middleware"
	"github.com/vlad-marlo/shortener/internal/store/filebased"
	"github.com/vlad-marlo/shortener/internal/store/inmemory"
	"github.com/vlad-marlo/shortener/internal/store/sqlstore"

	"github.com/vlad-marlo/shortener/internal/httpserver"
	"github.com/vlad-marlo/shortener/internal/store"
)

var (
	logLevel            logrus.Level = logrus.TraceLevel
	logOutput           io.Writer    = io.Discard
	logDir              string       = "logs"
	logDefaultFormatter string       = log.JSONFormatter
	logFormatter        *logrus.Formatter
)

func main() {
	storeLogger := createLogger("storage")
	serverLogger := createLogger("server")

	config := httpserver.NewConfig()

	storage, err := initStorage(config, storeLogger)
	if err != nil {
		serverLogger.Panicf("init storage: %v", err)
	}

	defer func() {
		if err := storage.Close(); err != nil {
			storeLogger.Panicf("close server: %v", err)
		}
	}()

	s := httpserver.New(config, storage, serverLogger)
	defer s.Close()
	serverLogger.WithFields(map[string]interface{}{
		"bind_addr": config.BindAddr,
	}).Info("successfully init server")

	go func() {
		// logging fatal because listen and server always return not-nil error
		serverLogger.Panicf("listen and server server: %v", s.ListenAndServe())
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)
	sig := <-interrupt
	serverLogger.WithFields(map[string]interface{}{
		"signal": sig.String(),
	}).Info("graceful shut down")
}

// createLogger creates new named logger with stdout and file output.
func createLogger(name string) *logrus.Logger {
	opts := []log.OptFunc{
		log.WithOutput(logOutput),
		log.WithLevel(logLevel),
		log.WithReportCaller(true),
		log.WithDefaultFormatter(logDefaultFormatter),
		log.WithHook(
			hook.New(
				logrus.AllLevels,
				[]io.Writer{os.Stdout},
				hook.WithFileOutput(
					logDir,
					name,
					time.Now().Format("2006-January-02-15"),
				),
			),
		),
	}
	if logFormatter != nil {
		opts = append(opts, log.WithFormatter(*logFormatter))
	}

	return log.WithOpts(opts...)
}

func initStorage(cfg *httpserver.Config, logger *logrus.Logger) (storage store.Store, err error) {
	switch cfg.StorageType {
	case store.InMemoryStorage:
		storage = inmemory.New()
	case store.FileBasedStorage:
		storage, err = filebased.New(cfg.FilePath)
	case store.SQLStore:
		storage, err = sqlstore.New(context.Background(), cfg.Database, logger)
	default:
		storage, err = filebased.New(cfg.FilePath)
	}
	return
}
