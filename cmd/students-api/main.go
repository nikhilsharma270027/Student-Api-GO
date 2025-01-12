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
	"github.com/nikhilsharma270027/GOLang-student-api/internal/http/handlers/student"
	"github.com/nikhilsharma270027/GOLang-student-api/internal/types/sqlite"
)

func main() {
	// fmt.Println("new")
	// load config
	cfg := config.MustLoad()

	// database setup
	// storage is database instance
	// , er := sqlite.New(cfg) changed after we need storage as prop
	storage, er := sqlite.New(cfg) //
	if er != nil {
		log.Fatal(er)
	}

	slog.Info("Storage Initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))
	// setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students", student.GetList(storage))
	router.HandleFunc("PUT /api/students/{id}", student.UpdateStudentById(storage))
	router.HandleFunc("DELETE /api/students/{id}", student.DeleteById(storage))

	// setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	// fmt.Printf("Envirnment started %s\n", cfg.Env)
	// fmt.Printf("Database started %s\n", cfg.StoragePath)
	slog.Info("Server started", slog.String("addressc", cfg.Addr))
	fmt.Printf("Server started on %s\n", cfg.Addr)

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	// its like if any interrupt by user or any other reason
	// send data in "done", it will stop the server

	// go func() {
	// 	err := server.ListenAndServe() // will start server
	// 	if err != nil {
	// 		log.Fatal("failed to start server", err)
	// 	}
	// }()
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server startup failed", slog.String("error", err.Error()))
			done <- os.Interrupt // Trigger shutdown if server can't start
		}
	}()

	<-done // untill done chan doesn't receiver any channel ,
	// we r stuck here , code wouldn't go further
	// the server will be running unless the done get some signal chan
	//
	// after done get data/trigger the next code will run

	slog.Info("Shuting down ")

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
