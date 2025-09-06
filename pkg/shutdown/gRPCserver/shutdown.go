package gRPCshutdown

import (
	"grpc/pkg/database"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Shutdown(db *database.Database) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop

	slog.Info("Starting graceful shutdown...", "signal", sign)

	done := make(chan bool, 1)
	go func() {
		if db != nil {
			db.Close()
		}
		done <- true
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
	}
	slog.Info("gRPC server stopped")
}
