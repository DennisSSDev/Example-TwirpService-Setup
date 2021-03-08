package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	server "github.com/dennisssdev/Example-TwirpService-Setup/internal"
	"github.com/dennisssdev/Example-TwirpService-Setup/rpc/example-service"
	"github.com/sirupsen/logrus"
)

func main() {
	twirpServer, err := server.NewExampleServer()
	if err != nil {
		logrus.Fatal("Failed to create example twirp server", err)
	}

	twirpHandler := example.NewExampleServer(twirpServer)
	mux := http.NewServeMux()
	mux.Handle(example.ExamplePathPrefix, twirpHandler)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, int64(10<<20)) // 10  MiB max per request
		mux.ServeHTTP(w, r)
	})

	port := 8000 // This should probably be in a config file
	appServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	shutdownCh := make(chan bool)
	go func() {
		shutdownSIGch := make(chan os.Signal, 1)
		signal.Notify(shutdownSIGch, syscall.SIGINT, syscall.SIGTERM)

		sig := <-shutdownSIGch

		logrus.Println("Received shutdown signal request", "signal", sig)
		if err := appServer.Shutdown(context.Background()); err != nil {
			logrus.Fatal("Failed to shutdown server", err)
		}
		close(shutdownCh)
	}()

	logrus.Println(fmt.Sprintf("Starting example twirp service on port %d", port))
	if err := appServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.Fatal("Failed to listen and serve", err)
	}

	<-shutdownCh
	logrus.Println("Server successfully shutdown")
}
