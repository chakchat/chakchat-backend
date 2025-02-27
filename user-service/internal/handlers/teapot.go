package handlers

import "net/http"

func AmITeapot(r *http.Request, w http.ResponseWriter) {
	w.WriteHeader(http.StatusTeapot)
	w.Write([]byte("I'm a teapot"))
}
