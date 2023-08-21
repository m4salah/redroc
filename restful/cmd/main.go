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
	downloadBackendAddress = flag.String("download_backend_address", "localhost:8080", "download backend server address")
	uploadBackendAddress   = flag.String("upload_backend_address", "localhost:8080", "upload backend server address")
	searchBackendAddress   = flag.String("search_backend_address", "localhost:8080", "search backend server address")
	backendTimeout         = flag.Duration("backend_timeout", 1*time.Minute, "backend request timeout")
	listenPort             = flag.Int("listen_port", 8080, "start server on this port")
	host                   = flag.String("host", "", "start server on this host")
	env                    = flag.String("env", "development", "Env")
	skiptGcloudAuth        = flag.Bool("skip_gcloud_auth", false, "disable gcloud auth")
)

// release is set through the linker at build time, generally from a git sha.
// Used for logging and error reporting.
var release string

func main() {
	os.Exit(start())
}

func start() int {
	flag.Parse()

	logger, err := util.CreateLogger(*env, release)
	if err != nil {
		log.Println("Error setting up the logger:", err)
		return 1
	}

	s := server.New(server.Options{
		SkipGcloudAuth:      *skiptGcloudAuth,
		Host:                *host,
		Log:                 logger,
		Port:                *listenPort,
		ConnTimeout:         *backendTimeout,
		DownloadBackendAddr: *downloadBackendAddress,
		UploadBackendAddr:   *uploadBackendAddress,
		SearchBackendAddr:   *searchBackendAddress,
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
