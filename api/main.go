package main

import (
	"api/handlers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind addres for the server")

func main() {

	// log object to pass on handlers struct
	l := log.New(os.Stdout, "api", log.LstdFlags)

	// create the handlers
	hh := handlers.NewHello(l)
	gh := handlers.NewGoodbye(l)

	// create a new serve mux and register the handler
	serveMux := http.NewServeMux()
	serveMux.Handle("/", hh)
	serveMux.Handle("/goodbye", gh)

	// create a new server
	server := &http.Server{
		Addr:         ":9090",           // configure the bind address
		Handler:      serveMux,          // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		IdleTimeout:  120 * time.Second, // max time to read request from the client
		ReadTimeout:  1 * time.Second,   // max time to write response to the client
		WriteTimeout: 1 * time.Second,   // max time for connections using TCP keep-Alive
	}

	// create goroutine for anonymous function to run concurrently
	// to not block gracefull shutdown
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			l.Fatal((err))
		}
	}()

	// passing the signal to notify close down
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// create channel
	sig := <-sigChan

	log.Println("Recieved terminated, graceful shutdown", sig)

	// 30 seconds to attempt to gracefully to shutdown, otherwise force close
	timeoutContext, _ := context.WithTimeout(context.Background(), 30*time.Second)

	// gracefully shutdown, awaits for connections to close
	// takes a context
	server.Shutdown(timeoutContext)
}
