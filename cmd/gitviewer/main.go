package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/theothertomelliott/gitviewer/pkg/server"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:\n gitviewer [repo path or URL]")
		return
	}
	repoPath := os.Args[1]

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
		log.Println("Serving on port :8080")
		err := http.ListenAndServe(":8080", nil)
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
