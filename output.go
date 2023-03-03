package main

import (
	"fmt"
	"net/http"
)

func writeMessage(w http.ResponseWriter, text string) {
	fmt.Fprintf(w, `{"code":200,"message":"%s"}`, text)
	w.WriteHeader(http.StatusOK)
	w.(http.Flusher).Flush()
}
func writeError(w http.ResponseWriter, err error) {
	fmt.Fprintf(w, `{"code":400,"error":"%s"}`, err.Error())
	w.WriteHeader(http.StatusBadRequest)
	w.(http.Flusher).Flush()
}
func writeDone(w http.ResponseWriter) {
	fmt.Fprintf(w, `{"code":200,"message":"done"}`)
	w.WriteHeader(http.StatusOK)
	w.(http.Flusher).Flush()
}
