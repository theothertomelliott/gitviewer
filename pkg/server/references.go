package server

import (
	"fmt"
	"net/http"

	"github.com/go-git/go-git/v5/plumbing"
)

func (s *Server) refs(w http.ResponseWriter, req *http.Request) {
	refs, _ := s.repo.References()
	refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() == plumbing.HashReference {
			fmt.Fprintf(w, "<a href=\"/revision/%v/\">%v</a> %v<br/>\n", ref.Hash(), ref.Hash(), ref.Name())
		}

		return nil
	})
}
