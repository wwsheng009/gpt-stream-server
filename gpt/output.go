package gpt

import (
	"net/http"
)

func writeMessage(w http.ResponseWriter, text string) {
	// fmt.Fprintf(w, `{"status":Success,"data":"%s"}`, text)
	w.Write([]byte(text))
	// w.WriteHeader(http.StatusOK)
	w.(http.Flusher).Flush()
}
func writeError(w http.ResponseWriter, err error) {
	// fmt.Fprintf(w, `{"status":Fail,"message":"%s"}`, err.Error())
	w.Write([]byte(err.Error()))
	w.WriteHeader(http.StatusBadRequest)
	w.(http.Flusher).Flush()
}
func writeDone(w http.ResponseWriter) {
	// fmt.Fprintf(w, `{"code":200,"message":"done"}`)
	w.Write([]byte("[DONE]"))
	w.WriteHeader(http.StatusOK)
	w.(http.Flusher).Flush()
}
