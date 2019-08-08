package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/douglasmakey/ursho/config"
	"github.com/douglasmakey/ursho/encoding"
	"github.com/douglasmakey/ursho/encoding/aes"
	"github.com/douglasmakey/ursho/encoding/base62"
	"github.com/douglasmakey/ursho/handler"
	"github.com/douglasmakey/ursho/storage"
	"github.com/douglasmakey/ursho/storage/dgraph"
	"github.com/douglasmakey/ursho/storage/postgres"
)

func main() {
	configPath := flag.String("config", "./config/config.json", "path of the config file")

	flag.Parse()

	// Read config
	config, err := config.FromFile(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	var svc storage.Service
	var coder encoding.Coder

	switch config.Options.Encoding {
	case "aes":
		key, err := hex.DecodeString(config.Crypto.Key)
		if err != nil {
			log.Fatal(err)
		}

		nonce, err := hex.DecodeString(config.Crypto.Nonce)
		if err != nil {
			log.Fatal(err)
		}

		c, err := aes.New(key, nonce)
		if err != nil {
			log.Fatal(err)
		}

		coder = c
	default:
		coder = base62.New()
	}

	switch config.Options.Database {
	case "dgraph":
		s, err := dgraph.New(config.Dgraph.Host, config.Dgraph.Port, coder)
		if err != nil {
			log.Fatal(err)
		}

		svc = s
	default:
		// Set use storage, select [Postgres, Filesystem, Redis ...]
		s, err := postgres.New(config.Postgres.Host, config.Postgres.Port, config.Postgres.User, config.Postgres.Password, config.Postgres.DB, coder)
		if err != nil {
			log.Fatal(err)
		}
		svc = s
	}

	defer svc.Close()

	// Create a server
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port),
		Handler: handler.New(config.Options.Prefix, svc),
	}

	go func() {
		// Start server
		log.Printf("Starting HTTP Server. Listening at %q", server.Addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("%v", err)
		} else {
			log.Println("Server closed!")
		}
	}()

	// Check for a closing signal
	// Graceful shutdown
	sigquit := make(chan os.Signal, 1)
	signal.Notify(sigquit, os.Interrupt, syscall.SIGTERM)

	sig := <-sigquit
	log.Printf("caught sig: %+v", sig)
	log.Printf("Gracefully shutting down server...")

	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("Unable to shut down server: %v", err)
	} else {
		log.Println("Server stopped")
	}
}
