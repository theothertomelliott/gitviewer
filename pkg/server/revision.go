package server

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func (s *Server) revision(w http.ResponseWriter, req *http.Request) {
	pathSegments := strings.SplitN(req.URL.Path, "/", 2)
	commitHash, filePath := req.URL.Path, ""
	if len(pathSegments) > 1 {
		commitHash, filePath = pathSegments[0], pathSegments[1]
	}

	commit, err := s.repo.CommitObject(plumbing.NewHash(commitHash))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// ... retrieve the tree from the commit
	tree, err := commit.Tree()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(filePath) > 0 && !strings.HasSuffix(filePath, "/") {
		err = s.serveFileContent(w, tree, filePath)
	} else {
		err = s.serveDirectoryListing(w, tree, filePath)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) serveFileContent(w http.ResponseWriter, tree *object.Tree, filePath string) error {
	file, err := tree.File(filePath)
	if err != nil {
		return err
	}
	contents, err := file.Contents()
	if err != nil {
		return err
	}
	fmt.Fprint(w, contents)
	return nil
}

func (s *Server) serveDirectoryListing(w http.ResponseWriter, tree *object.Tree, filePath string) error {
	fileList := make(map[string]struct{})
	dirList := make(map[string]struct{})
	tree.Files().ForEach(func(f *object.File) error {
		if !strings.HasPrefix(f.Name, filePath) {
			return nil
		}
		childPath := strings.Replace(f.Name, filePath, "", 1)
		if !strings.Contains(childPath, "/") {
			fileList[childPath] = struct{}{}
		} else {
			pathDir := strings.Split(childPath, "/")[0]
			dirList[fmt.Sprintf("%v/", pathDir)] = struct{}{}
		}
		return nil
	})

	for _, dir := range sortSet(dirList) {
		fmt.Fprintf(w, "<a href=\"%v\">%v</a><br/>", dir, dir)
	}
	for _, file := range sortSet(fileList) {
		fmt.Fprintf(w, "<a href=\"%v\">%v</a><br/>", file, file)
	}
	return nil
}

func sortSet(set map[string]struct{}) []string {
	var out []string
	for s := range set {
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}
