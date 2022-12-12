package main

import (
	"context"
	"fmt"
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
	logLevel            = logrus.TraceLevel
	logOutput           = io.Discard
	logDir              = "logs"
	logDefaultFormatter = log.JSONFormatter
	logFormatter        *logrus.Formatter
	buildVersion        = "N/A"
	buildDate           = "N/A"
	buildCommit         = "N/A"
	logFileNameFormat   = "2006-January-02-15"
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
	debugInfo()

	s := httpserver.New(config, storage, serverLogger)
	defer func() {
		if err := s.Close(); err != nil {
			serverLogger.Fatal(fmt.Sprintf("close server: %v", err))
		}
	}()
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
func createLogger(name string) *logrus.Entry {
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
					time.Now().Format(logFileNameFormat),
				),
			),
		),
	}
	if logFormatter != nil {
		opts = append(opts, log.WithFormatter(*logFormatter))
	}

	return log.WithOpts(opts...).WithFields(logrus.Fields{
		"build_version": buildVersion,
		"build_commit":  buildCommit,
		"build_date":    buildDate,
	})
}

// initStorage is abstract factory to create new storage object with provided config.
func initStorage(cfg *httpserver.Config, logger *logrus.Entry) (storage store.Store, err error) {
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

// debugInfo ...
func debugInfo() {
	fmt.Printf("Build version: %s \nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)
}
