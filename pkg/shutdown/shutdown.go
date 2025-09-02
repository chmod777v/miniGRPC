package shutdown

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Shutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop

	log.Println("Service terminated signal:", sign)

}
