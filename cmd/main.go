package main

import (
	"fmt"
	"github.com/Fyefhqdishka/deadlock_v.2/internal/app"
	"github.com/Fyefhqdishka/deadlock_v.2/internal/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("can't load config, err: %v", err)
	}

	app, err := app.New(cfg)
	if err != nil {
		log.Fatalf("can't load server, err: %v", err)
	}

	go func() {
		if err = app.Run(); err != nil {
			log.Fatalf("server failed: %v", err)
		}
	}()
	log.Println("shutting down...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	if err := app.Stop(); err != nil {
		fmt.Errorf("error during shutdown: %v", err)
	}
}
