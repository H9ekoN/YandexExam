package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/H9ekoN/YandexExam_go/internal/application"
)

func main() {
	port := flag.String("port", "", "Server port (default: from PORT env var or 8080)")
	flag.Parse()

	cfg := &application.Config{
		Port: *port,
	}
	if cfg.Port == "" {
		cfg.Port = os.Getenv("PORT")
	}

	app := application.NewServer(cfg)

	done := make(chan bool, 1)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals
		fmt.Printf("\nReceived signal: %v\n", sig)
		done <- true
	}()

	go func() {
		log.Printf("Starting calculator service on port %s\n", app.Port())
		if err := app.Start(); err != nil {
			log.Printf("Server error: %v\n", err)
			done <- true
		}
	}()

	<-done

	log.Println("Shutting down...")
	if err := app.Stop(); err != nil {
		log.Printf("Error during shutdown: %v\n", err)
		os.Exit(1)
	}
	log.Println("Server stopped successfully")
}
