package main

import (
	"fmt"
	"path/filepath"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var (
	mainFileDir = "cmd/api/"
	handlerDir  = "handler/"
	serverDir   = "server/"
	viewsDir    = "views/"

	mainFilePath    = filepath.Join(mainFileDir, "main.go")
	handlerFilePath = filepath.Join(handlerDir, "handler.go")
	serverFilePath  = filepath.Join(serverDir, "server.go")
	routesFilePath  = filepath.Join(serverDir, "routes.go")
	makeFilePath    = "Makefile"
	envFilePath     = ".env"
	goModFilePath   = "go.mod"
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
			log.Println("Server forced to shutdown with error:", err)
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
				panic(fmt.Sprintln("http server error:", err))
			}

			<-done

			log.Println("Graceful shutdown complete.")
		}
	`, name)
	result["main.go"] = mainFileData

	// TODO: Do the same for the other files in `template-files/`

	return result, nil
}

func main() {
	data, err := addNameToFiles("foobar")
	check(err)
	fmt.Println(data)

	// TODO:
	// 1. Accept flags to add the project name and choose custom options.
	// 2. Create all the necessary directories and files and write the updated data to them.
}
