package main

import (
	"context"
	"fmt"
	"ilkerciblak/socketoid/internal/api"
	"ilkerciblak/socketoid/internal/platform"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	errchan := make(chan error, 1)
	app_config, _ := platform.Config()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	server, _ := api.Server(
		app_config.ADDR,
		app_config.IdleTimeout,
	)

	server.Start(errchan)

	select {
	case sign := <-sig:
		fmt.Println("signal received: ", sign)
		close(sig)
		server.Shutdown(ctx)
		return
	case err := <-errchan:
		fmt.Println("server failed due: ", err)
		return
	}
}
