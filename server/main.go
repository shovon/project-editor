package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"server/config"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func createFile(filePath string) (*os.File, error) {
	dir := filepath.Dir(filePath)

	// Check if the directory already exists
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		// Create the directory and any necessary parent directories
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory: %s", err)
		}
	}

	return os.Create(filePath)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

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
				// TODO: this should really be thrown into a job queue.
				file, err := createFile(path.Join(config.FolderPath(), r.URL.Path))
				if err != nil {
					http.Error(w, "Error creating file.", http.StatusInternalServerError)
					return
				}
				defer file.Close()

				_, err = io.Copy(file, r.Body)
				if err != nil {
					http.Error(w, "Error streaming data to file.", http.StatusInternalServerError)
					return
				}

				w.WriteHeader(200)
				w.Write([]byte("{ \"status\": \"ok\" }"))
			}),
		).ServeHTTP(w, r)
	})

	r.Delete("/files/*", func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		http.StripPrefix(
			pathPrefix,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				err := os.Remove(path.Join(config.FolderPath(), r.URL.Path))
				if err != nil {
					http.Error(w, "Error deleting file.", http.StatusInternalServerError)
					return
				}

				w.WriteHeader(200)
				w.Write([]byte("{ \"status\": \"ok\" }"))
			}),
		).ServeHTTP(w, r)
	})

	port := config.Port()
	fmt.Printf("Server listening on port :%d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
