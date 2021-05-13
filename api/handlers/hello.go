package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// simple handler
type Hello struct {
	l *log.Logger
}

// created a new hello handler with the given logger
func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

// Define method in a struct
// ServeHTTP implements the go http.handler interface
func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	log.Println("Running Hello Handler")

	// read the body
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading body", err)

		http.Error(rw, "Unable to read request body", http.StatusBadRequest)
		return
	}

	// write the response
	fmt.Fprintf(rw, "Hello %s", b)
}
