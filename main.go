package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"

	_ "github.com/aakordas/pservice/ptlist"
)

func main() {
	var port string

	args := os.Args[1:]
	nArgs := len(args)

	logger := log.New(os.Stderr, "",
		log.Ldate|log.Ltime|log.Lshortfile)

	// Not really sure if/how I could utilise that with Docker, no
	// time to find out, either.
	if nArgs == 0 {
		port = "localhost:8282"
	} else if nArgs == 1 {
		port = args[0]
	} else {
		logger.Println("Provide an addr/host argument, only.")
		logger.Fatalln("If no arguments are provided, the default is \"localhost:8282\"")
	}

	matched, err := regexp.MatchString(`.*:\d{4}`, port)
	if err != nil {
		logger.Fatalln(err)
	}

	if !matched {
		logger.Fatalln("The first argument must have the form addr:port")
	}

	server := &http.Server{
		Addr:         port[strings.Index(port, `:`):],
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, os.Kill)
		<-sigint

		if err := server.Shutdown(context.Background()); err != nil {
			log.Println("HTTP server Shutdown: %v", err)
		}
	}()

	log.Println("Server listening at " + port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalln("HTTP server ListenAndServe: %v", err)
	} else {
		log.Println("HTTP server shut down successfully.")
	}
}
