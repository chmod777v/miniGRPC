package httpshutdown

import (
	grpcconect "grpc/pkg/httpServer/grpcConect"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Shutdown(serv *grpcconect.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop

	slog.Info("Starting graceful shutdown...", "signal", sign)

	done := make(chan bool, 1)
	go func() {
		serv.Close()
		done <- true
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
	}
	slog.Info("Http server stopped")
}
