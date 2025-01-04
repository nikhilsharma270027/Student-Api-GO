package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nikhilsharma270027/GOLang-student-api/internal/config"
)

func main() {
	// fmt.Println("new")
	// load config
	cfg := config.MustLoad()

	// database setup
	// setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to student api"))
	})

	// setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	// fmt.Printf("Envirnment started %s\n", cfg.Env)
	// fmt.Printf("Database started %s\n", cfg.StoragePath)
	slog.Info("Server started", slog.String("addressc", cfg.Addr))
	fmt.Printf("Server started %s\n", cfg.HTTPServer)

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	// its like if any interrupt by user or any other reason
	// send data in "done", it will stop the server

	go func() {
		err := server.ListenAndServe() // will start server
		if err != nil {
			log.Fatal("failed to start server", err)
		}
	}()

	<-done // untill done chan doesn't receiver any channel ,
	// we r stuck here , code wouldn't go further
	// the server will be running unless the done get some signal chan
	//
	// after done get data/trigger the next code will run

	slog.Info("Shuuting down ")

	// we use timer if after specific time shtdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// ctx is context, in return we get cancel method()
	defer cancel()

	err := server.Shutdown(ctx) // we have code the shutddown the server
	// but can infinitly block port for a while maybe
	if err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	//---- Same code as above
	// if err := server.Shutdown(ctx);  err != nil {
	// 	slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	// }

	slog.Info("Server shtdown successfully")
}
