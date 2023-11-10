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
	"github.com/m4salah/redroc/apps/server/types"
	"github.com/m4salah/redroc/libs/util"
	"golang.org/x/sync/errgroup"
)

var (
	backendTimeout  = flag.Duration("backend_timeout", 1*time.Minute, "backend request timeout")
	env             = flag.String("env", util.LOCALENV, "Env")
	skiptGcloudAuth = flag.Bool("skip_gcloud_auth", false, "disable gcloud auth")
)

// release is set through the linker at build time, generally from a git sha.
// Used for logging and error reporting.
var release string

var config types.Config

func main() {
	if err := util.LoadConfig(&config); err != nil {
		panic(err)
	}
	os.Exit(start())
}

func start() int {
	flag.Parse()

	util.InitializeSlog(*env, release)
	s := server.New(server.Options{
		SkipGcloudAuth: *skiptGcloudAuth,
		ConnTimeout:    *backendTimeout,
		ServerConfig:   config,
	})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := s.Start(); err != nil {
			slog.Error("Error starting server", "error", err)
			return err
		}
		return nil
	})

	<-ctx.Done()

	eg.Go(func() error {
		if err := s.Stop(); err != nil {
			slog.Error("Error stopping server", "error", err)
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return 1
	}
	return 0
}
