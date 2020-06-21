package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/theothertomelliott/gitviewer/pkg/server"
)

var port = flag.Int("port", 8080, "port to listen on")

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage:\n gitviewer [repo path or URL]")
		return
	}
	repoPath := args[0]

	var (
		s   *server.Server
		err error
	)

	if isValidURL(repoPath) {
		s, err = server.NewFromURL(repoPath)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		s, err = server.NewFromRepoFolder(repoPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	defer func() {
		err := s.Shutdown()
		if err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	h, err := s.GetHandler()
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", h)

	go func() {
		log.Printf("Serving on port :%d", *port)
		err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
	<-c
}

// isValidURL tests a string to determine if it is a well-structured url or not.
func isValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}
