package server

import (
	"fmt"
	"net/http"
)

func (s *Server) index(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, `
		<a href="/references">References</a><br/>
		<a href="/commits">Commits</a><br/>
	`)
}
