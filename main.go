package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var (
	mainFileDir = "test/cmd/api/"
	handlerDir  = "test/handler/"
	serverDir   = "test/server/"
	viewsDir    = "test/views/"
	// mainFileDir = "cmd/api/"
	// handlerDir  = "handler/"
	// serverDir   = "server/"
	// viewsDir    = "views/"

	mainFilePath      = filepath.Join(mainFileDir, "main.go")
	handlerFilePath   = filepath.Join(handlerDir, "handler.go")
	serverFilePath    = filepath.Join(serverDir, "server.go")
	routesFilePath    = filepath.Join(serverDir, "routes.go")
	viewsBaseFilePath = filepath.Join(viewsDir, "base.templ")
	viewsHomeFilePath = filepath.Join(viewsDir, "home.templ")
	makeFilePath      = "test/Makefile"
	envFilePath       = "test/.env"
	goModFilePath     = "test/go.mod"
	// makeFilePath      = "Makefile"
	// envFilePath       = ".env"
	// goModFilePath     = "go.mod"
)

func addNameToFiles(name string) (map[string]string, error) {
	result := make(map[string]string, 50)
	mainFileData := fmt.Sprintf(`
		package main

		import (
			"context"
			"fmt"
			"log"
			"net/http"
			"os/signal"
			"syscall"
			"time"

			"github.com/devkaare/%s/server"
		)

		func gracefulShutdown(apiServer *http.Server, done chan bool) {
			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			<-ctx.Done()

			log.Println("shutting down gracefully, press Ctrl+C again to force")

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := apiServer.Shutdown(ctx); err != nil {
				log.Printf("Server forced to shutdown with error :%%v", err)
			}

			log.Println("Server exiting")

			done <- true
		}

		func main() {
			server := server.NewServer()

			done := make(chan bool, 1)

			go gracefulShutdown(server, done)

			err := server.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				panic(fmt.Sprintf("http server error: %%s", err))
			}

			<-done

			log.Println("Graceful shutdown complete.")
		}
	`, name)
	result[mainFilePath] = mainFileData

	handlerFileData := fmt.Sprintf(`
		package handler

		import (
			"net/http"

			"github.com/devkaare/%s/farms"
		)

		type New struct{}

		func (t *New) Greet(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello World!"))
		}
	`, name)
	result[handlerFilePath] = handlerFileData

	serverFileData := fmt.Sprintf(`
		package server

		import (
			"fmt"
			"net/http"
			"os"
			"strconv"
			"time"

			_ "github.com/joho/godotenv/autoload"
		)

		type Server struct {
			port int
		}

		func NewServer() *http.Server {
			port, _ := strconv.Atoi(os.Getenv("PORT"))
			NewServer := &Server{
				port: port,
			}

			// Declare Server config
			server := &http.Server{
				Addr:         fmt.Sprintf(":%%d", NewServer.port),
				Handler:      NewServer.RegisterRoutes(),
				IdleTimeout:  10 * time.Second,
				WriteTimeout: 30 * time.Second,
			}

			return server
		}
	`)
	result[serverFilePath] = serverFileData

	routesFileData := fmt.Sprintf(`
		package server

		import (
			"net/http"

			"github.com/devkaare/%s/handler"
			"github.com/go-chi/chi/v5"
			"github.com/go-chi/chi/v5/middleware"
			"github.com/go-chi/cors"
		)

		func (s *Server) RegisterRoutes() http.Handler {
			r := chi.NewRouter()
			r.Use(middleware.Logger)

			r.Use(cors.Handler(cors.Options{
				AllowedOrigins:   []string{"https://*", "http://*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
				AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
				AllowCredentials: true,
				MaxAge:           300,
			}))

			r.Route("/", s.RegisterNewRoutes)

			return r
		}

		func (s *Server) RegisterNewRoutes(r chi.Router) {
			handler := &handler.New{}

			r.Get("/hello", handler.Greet)

		}
	`, name)
	result[routesFilePath] = routesFileData

	result[viewsBaseFilePath] = `
		package views

		templ Base() {
			<!DOCTYPE html>
			<html lang="en" >
				<head>
					<meta charset="utf-8"/>
					<meta name="viewport" content="width=device-width,initial-scale=1"/>
					<title>Welcome</title>
				</head>
				<body >
					<main >
						{ children... }
					</main>
				</body>
			</html>
		}
	`

	result[viewsHomeFilePath] = `
		package views

		templ Home() {
			@Base() {
				</p>Hello World!</p>
			}
		}
	`

	makeFileData := fmt.Sprintf(`
		MAIN_FILE_PATH = %s

		all: build test

		run:
			@go run $(MAIN_FILE_PATH)

		build:
			@echo "Building..."
			@go build $(MAIN_FILE_PATH)

		test:
			@echo "Testing..."
			@go test ./... -v

		clean:
			@echo "Cleaning..."
			@rm -rf main
			@go mod tidy

		.PHONY: all run build test clean
	`, mainFilePath)
	result[makeFilePath] = makeFileData

	result[envFilePath] = `
		PORT=8080
	`

	result[goModFilePath] = fmt.Sprintf(`
		module github.com/devkaare/%s

		go 1.23.5

		require (
			github.com/a-h/templ v0.3.833
			github.com/go-chi/chi/v5 v5.2.1
			github.com/go-chi/cors v1.2.1
			github.com/joho/godotenv v1.5.1
		)
	`, name)

	return result, nil
}

func main() {
	rawData, err := addNameToFiles("foobar")
	check(err)

	fmt.Println("Creating folders") // Create folders
	os.MkdirAll(mainFileDir, 0755)
	os.Mkdir(handlerDir, 0755)
	os.Mkdir(serverDir, 0755)
	os.Mkdir(viewsDir, 0755)

	fmt.Println("Creating and writing files") // Create files and write data
	for filePath, data := range rawData {
		fmt.Println("File path:", filePath, data)
		os.WriteFile(filePath, []byte(data), 0644)
	}
}

