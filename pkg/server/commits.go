package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"
)

func (s *Server) commits(w http.ResponseWriter, req *http.Request) {
	commits, _ := s.repo.CommitObjects()
	commits.ForEach(func(commit *object.Commit) error {
		fmt.Fprintf(w, "<a href=\"/revision/%v/\">%v</a> %v<br/>\n", commit.Hash.String(), commit.Hash.String(), strings.Split(commit.Message, "\n")[0])
		return nil
	})
}
