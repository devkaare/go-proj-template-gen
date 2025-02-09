package handler

import "net/http"

type New struct{}

func (t *New) Greet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}