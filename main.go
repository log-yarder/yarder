package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/document"
)

func main() {
	tmpDir, err := ioutil.TempDir("", "yarder-dev")
	if err != nil {
		log.Panicf(fmt.Sprintf("Unable to create temp dir, %v", err))
	}

	indexer, err := bleve.New(tmpDir, bleve.NewIndexMapping())
	if err != nil {
		log.Panicf("Unable to intialise bleve: %v", err)
	}

	handler := &handler{indexer: indexer}
	srv := httptest.NewServer(handler)
	defer srv.Close()

	http.Post(srv.URL, "application/json", strings.NewReader(`{"message": "a test entry"}`))
	resp, _ := http.Get(srv.URL + "?q=entry")
	io.Copy(os.Stdout, resp.Body)
}

type handler struct {
	indexer bleve.Index
	lastID  int
}

func (s *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if strings.TrimRight(req.URL.Path, "/") != "" {
		http.NotFound(w, req)
		return
	}
	if req.Method == "GET" {
		s.queryHandler(w, req)
		return
	}
	if req.Method == "POST" {
		s.postHandler(w, req)
		return
	}
	http.NotFound(w, req)
}

func (s *handler) postHandler(w http.ResponseWriter, req *http.Request) {
	dec := json.NewDecoder(req.Body)
	var data interface{}
	err := dec.Decode(&data)
	req.Body.Close()
	if err != nil {
		s.writeError(w, err)
		return
	}

	doc := document.NewDocument(strconv.Itoa(s.lastID))
	s.lastID++
	err = s.indexer.Mapping().MapDocument(doc, data)
	if err != nil {
		s.writeError(w, err)
		return
	}

	index, _, err := s.indexer.Advanced()
	if err != nil {
		s.writeError(w, err)
		return
	}

	err = index.Update(doc)
	if err != nil {
		s.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *handler) queryHandler(w http.ResponseWriter, req *http.Request) {
	search := bleve.NewSearchRequest(bleve.NewQueryStringQuery(req.Header.Get("q")))
	results, err := s.indexer.Search(search)
	if err != nil {
		s.writeError(w, err)
	}

	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	err = enc.Encode(results)
	if err != nil {
		// can no longer serve error response
		log.Printf("ERROR: %v", err)
	}
}

func (s *handler) writeError(w http.ResponseWriter, err error) {
	log.Printf("ERROR: %v", err)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}
