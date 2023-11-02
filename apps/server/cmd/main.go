package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/m4salah/redroc/apps/server/server"
	"github.com/m4salah/redroc/libs/util"
	"golang.org/x/sync/errgroup"
)

var (
	backendTimeout  = flag.Duration("backend_timeout", 1*time.Minute, "backend request timeout")
	listenPort      = flag.Int("listen_port", 8080, "start server on this port")
	host            = flag.String("host", "", "start server on this host")
	env             = flag.String("env", "development", "Env")
	skiptGcloudAuth = flag.Bool("skip_gcloud_auth", false, "disable gcloud auth")
)

// release is set through the linker at build time, generally from a git sha.
// Used for logging and error reporting.
var release string

var config Config

type Config struct {
	DownloadBackendAddr string `env:"DOWNLOAD_BACKEND_ADDR,notEmpty"`
	UploadBackendAddr   string `env:"UPLOAD_BACKEND_ADDR,notEmpty"`
	SearchBackendAddr   string `env:"SEARCH_BACKEND_ADDR,notEmpty"`
}

func main() {
	err := util.LoadConfig(&config)
	if err != nil {
		panic(err)
	}
	os.Exit(start())
}

func start() int {
	flag.Parse()

	util.InitializeSlog(*env, release)
	s := server.New(server.Options{
		SkipGcloudAuth:      *skiptGcloudAuth,
		Host:                *host,
		Port:                *listenPort,
		ConnTimeout:         *backendTimeout,
		DownloadBackendAddr: config.DownloadBackendAddr,
		UploadBackendAddr:   config.UploadBackendAddr,
		SearchBackendAddr:   config.SearchBackendAddr,
	})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := s.Start(); err != nil {
			slog.Error("Error starting server", slog.String("err", err.Error()))
			return err
		}
		return nil
	})

	<-ctx.Done()

	eg.Go(func() error {
		if err := s.Stop(); err != nil {
			slog.Error("Error stopping server", slog.String("err", err.Error()))
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return 1
	}
	return 0
}
