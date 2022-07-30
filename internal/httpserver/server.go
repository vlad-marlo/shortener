package httpserver

import "net/http"

func Start(addr string) error {
	http.HandleFunc("/", handleUrlGetCreate)
	return http.ListenAndServe(addr, nil)
}
