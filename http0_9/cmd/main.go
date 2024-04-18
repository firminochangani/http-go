package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"golang.org/x/sync/errgroup"
	"http_v0_9/server"
)

type Config struct {
	Port int
	Host string
	Root string
}

func main() {
	config := &Config{}
	logLevel := flag.String("log-level", "info", "Log level")
	flag.IntVar(&config.Port, "port", 80, "Port number")
	flag.StringVar(&config.Host, "host", "127.0.0.1", "Host")
	flag.StringVar(&config.Root, "root", "./", "Root root")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: mapToLogLevel(*logLevel),
	}))

	logger.Info("initialising with the following options", "options", config)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()

	if err := run(ctx, config, logger); err != nil {
		logger.Error("shutting down server", "error", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func run(ctx context.Context, config *Config, logger *slog.Logger) error {
	g, gCtx := errgroup.WithContext(ctx)

	httpServer, err := server.NewServer(config.Host, config.Port, config.Root, logger)
	if err != nil {
		return err
	}

	g.Go(func() error {
		logger.Info("the http server is running")
		if err := httpServer.Listen(); !errors.Is(err, server.ErrHttpServerClosed) {
			return err
		}

		return nil
	})

	g.Go(func() error {
		<-gCtx.Done()
		logger.Info("terminating server gracefully")
		return httpServer.Shutdown()
	})

	return g.Wait()
}

func mapToLogLevel(value string) slog.Level {
	switch strings.ToLower(value) {
	case "debug":
		return slog.LevelDebug
	default:
		return slog.LevelInfo
	}
}
