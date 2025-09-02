package shutdown

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func Shutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop

	slog.Info("Service terminated", "signal", sign)

}
