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

	"github.com/vlad-marlo/shortener/internal/httpserver"
	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/filebased"
	"github.com/vlad-marlo/shortener/internal/store/inmemory"
	"github.com/vlad-marlo/shortener/internal/store/sqlstore"
)

func main() {
	serverLogger := log.WithOpts(
		log.WithOutput(io.Discard),
		log.WithLevel(logrus.TraceLevel),
		log.WithReportCaller(true),
		log.WithDefaultFormatter(log.JSONFormatter),
		log.WithHook(
			hook.New(
				logrus.AllLevels,
				[]io.Writer{os.Stdout},
				hook.WithFileOutput(
					"logs",
					"server",
					time.Now().Format("2006-January-02-15"),
				),
			),
		),
	)

	storeLogger := log.WithOpts(
		log.WithOutput(io.Discard),
		log.WithLevel(logrus.TraceLevel),
		log.WithReportCaller(true),
		log.WithDefaultFormatter(log.JSONFormatter),
		log.WithHook(
			hook.New(
				logrus.AllLevels,
				[]io.Writer{os.Stdout},
				hook.WithFileOutput(
					"logs",
					"storage",
					time.Now().Format("2006-January-02-15"),
				),
			),
		),
	)

	config := httpserver.NewConfig()

	var storage store.Store
	var err error

	switch config.StorageType {
	case store.InMemoryStorage:
		storage, err = inmemory.New(), nil
	case store.FileBasedStorage:
		storage, err = filebased.New(config.FilePath)
	case store.SQLStore:
		storage, err = sqlstore.New(context.Background(), config.Database, storeLogger)
	default:
		storage, err = filebased.New(config.FilePath)
	}
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
