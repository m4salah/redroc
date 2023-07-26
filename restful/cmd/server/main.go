package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/m4salah/redroc/restful/server"
	"github.com/m4salah/redroc/util"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var (
	backendAddress = flag.String("backend_address", "localhost:8080", "backend server address")
	backendTimeout = flag.Duration("backend_timeout", 10*time.Second, "backend request timeout")
	listenPort     = flag.Int("listen_port", 3300, "start server on this port")
	host           = flag.String("host", "localhost", "start server on this host")
	env            = flag.String("env", "development", "Env")
)

// release is set through the linker at build time, generally from a git sha.
// Used for logging and error reporting.
var release string

func main() {
	os.Exit(start())
}

func start() int {
	flag.Parse()

	logger, err := util.CreateLogger(*env)
	if err != nil {
		log.Println("Error setting up the logger:", err)
		return 1
	}
	logger = logger.With(zap.String("release", release))
	s := server.New(server.Options{
		Host:                *host,
		Log:                 logger,
		Port:                *listenPort,
		ConnTimeout:         *backendTimeout,
		DownloadBackendAddr: *backendAddress,
	})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := s.Start(); err != nil {
			logger.Info("Error starting server", zap.Error(err))
			return err
		}
		return nil
	})

	<-ctx.Done()

	eg.Go(func() error {
		if err := s.Stop(); err != nil {
			logger.Info("Error stopping server", zap.Error(err))
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return 1
	}
	return 0
}
