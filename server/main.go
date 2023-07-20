package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/config"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/file-paths", func(w http.ResponseWriter, r *http.Request) {
		filenames, err := getAllFilenames(config.FolderPath())
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Print all filenames with their full paths
		for _, filename := range filenames {
			fmt.Println(filename)
		}

		result, err := json.Marshal(filenames)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	})
	filesDir := http.Dir(config.FolderPath())

	r.Get("/files/*", func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		http.StripPrefix(pathPrefix, http.FileServer(filesDir)).ServeHTTP(w, r)
	})

	r.Put("/files/*", func(w http.ResponseWriter, r *http.Request) {
		// Here, we take the route, and upsert the incoming data from the request
		// body into the file at the path specified by the route.

		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		http.StripPrefix(
			pathPrefix,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Println(r.URL.Path)
			}),
		)
	})

	port := config.Port()
	fmt.Printf("Server listening on port :%d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
