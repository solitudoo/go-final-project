package server

import (
	"go-final-project/pkg/api"
	"net/http"
	"os"
)

func Run() error {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	addr := ":" + port

	webDir := "./web"
	handler := http.FileServer(http.Dir(webDir))

	http.Handle("/", handler)
	api.Init()

	return http.ListenAndServe(addr, nil)
}
