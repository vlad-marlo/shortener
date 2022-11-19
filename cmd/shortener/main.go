package main

import (
	"github.com/sirupsen/logrus"
	log "github.com/vlad-marlo/logger"
	"github.com/vlad-marlo/logger/hook"
	"github.com/vlad-marlo/shortener/internal/httpserver"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	serverLogger := log.WithOpts(
		log.WithOutput(io.Discard),
		log.WithLevel(logrus.DebugLevel),
		log.WithReportCaller(true),
		log.WithDefaultFormatter(log.JSONFormatter),
		log.WithHook(
			hook.New(
				logrus.AllLevels,
				nil,
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
		log.WithLevel(logrus.DebugLevel),
		log.WithReportCaller(true),
		log.WithDefaultFormatter(log.JSONFormatter),
		log.WithHook(
			hook.New(
				logrus.AllLevels,
				nil,
				hook.WithFileOutput(
					"logs",
					"server",
					time.Now().Format("2006-January-02-15"),
				),
			),
		),
	)

	config := httpserver.NewConfig()
	go func() {
		if err := httpserver.Start(config, serverLogger, storeLogger); err != nil {
			serverLogger.WithFields(map[string]interface{}{
				"bind_addr": config.BindAddr,
			}).Fatal(err)
		}
	}()
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	sig := <-interrupt
	serverLogger.WithFields(map[string]interface{}{
		"signal": sig.String(),
	}).Info("graceful shut down")
}
