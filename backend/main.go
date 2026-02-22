package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/apaul13/manual-labor/api"
	"github.com/apaul13/manual-labor/database"
)

func main() {
	// Initialize the global DB pool. This must succeed before starting the server.
	if err := database.InitDB(); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	// Ensure we close the DB pool when the process exits.
	defer func() {
		if err := database.CloseDB(); err != nil {
			log.Printf("error closing database: %v", err)
		}
	}()

	// Start the router in a goroutine so we can listen for shutdown signals here.
	go api.RunRouter()

	// Wait for interrupt/terminate signal.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	log.Printf("received signal %v, shutting down", s)

	// Deferred CloseDB will run; give a short log and exit.
	log.Println("shutdown complete")
}
