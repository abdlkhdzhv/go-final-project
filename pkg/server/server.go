package server

import "net/http"

func Run() error {
	webDir := "./web"
	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	return http.ListenAndServe(":7540", nil)
}
