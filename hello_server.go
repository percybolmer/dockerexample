package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	if name == "" {
		name = "Guest"
	}

	cmd := exec.Command("python", "run_this.py")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	// ExeC FfMPeg
	cmd = exec.Command("ffmpeg")
	out, err = cmd.CombinedOutput()
	if err != nil {
		log.Printf("Err: %s", err.Error())
	}

	log.Printf("Recieved request for %s\n", name)
	w.Write([]byte(fmt.Sprintf("%s, %s\n", string(out), name)))
}

func main() {
	// Create server and Route handlers
	r := mux.NewRouter()

	r.HandleFunc("/", handler)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Logging
	LOG_FILE_LOCATION := os.Getenv("LOG_FILE_LOCATION")
	if LOG_FILE_LOCATION != "" {
		log.SetOutput(&lumberjack.Logger{
			Filename:   LOG_FILE_LOCATION,
			MaxSize:    500,
			MaxBackups: 3,
			MaxAge:     28,
			Compress:   true,
		})
	}
	// Start server
	go func() {
		log.Println("Starting Server")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Handle Graceful Shutdown
	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	// this is Cool, We create a channel taking OS signals
	// And then we apply a Notifiaction to those signals
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block channel since its only 1 item
	<-interruptChan

	// Create deadline and wait for it
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting Down gracaefully")
	os.Exit(0)
}
