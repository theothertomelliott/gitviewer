package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-git/go-git/v5"
)

// Server handles serving of a Git Repository via HTTP
type Server struct {
	repo     *git.Repository
	shutdown func() error
}

// NewFromRepoFolder creates a server for a local folder containing a Git repo
func NewFromRepoFolder(repoPath string) (*Server, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}
	return &Server{
		repo: repo,
	}, nil
}

// NewFromURL creates a server for a given Git URL.
// The repo is cloned out to a temporary folder.
// NewFromURL will block until the repo has been cloned.
func NewFromURL(repoURL string) (*Server, error) {
	// Tempdir to clone the repository
	dir, err := ioutil.TempDir("", "clone-example")
	if err != nil {
		return nil, err
	}

	log.Printf("Cloning repo from %q into %q", repoURL, dir)

	cleanup := func() error {
		log.Printf("Cleaning up %q", dir)
		return os.RemoveAll(dir) // clean up
	}

	// Clones the repository into the given dir, just as a normal git clone does
	repo, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL: repoURL,
	})

	if err != nil {
		cleanup()
		return nil, err
	}

	return &Server{
		repo:     repo,
		shutdown: cleanup,
	}, nil
}

// GetHandler creates an http.Handler to serve requests for
// content in the server's repo.
func (s *Server) GetHandler() (http.Handler, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.index)
	mux.Handle("/revision/", http.StripPrefix("/revision/", http.HandlerFunc(s.revision)))
	mux.HandleFunc("/references", s.refs)
	mux.HandleFunc("/commits", s.commits)
	return mux, nil
}

// Shutdown cleans up any temporary resources used.
func (s *Server) Shutdown() error {
	if s.shutdown != nil {
		return s.shutdown()
	}
	return nil
}
