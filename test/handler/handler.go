
		package handler

		import (
			"net/http"

			"github.com/devkaare/foobar/farms"
		)

		type New struct{}

		func (t *New) Greet(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello World!"))
		}
	