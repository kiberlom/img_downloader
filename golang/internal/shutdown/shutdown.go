package shutdown

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func NewShutdown() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		done := make(chan os.Signal, 1)
		signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
		<-done
		log.Println("Завершение работы")
		cancel()
	}()

	return ctx

}
