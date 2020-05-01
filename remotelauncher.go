package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func Shutdown(writer http.ResponseWriter, req *http.Request) {
	log.Printf("Got shutdown request")
	when, ok := req.URL.Query()["when"]
	if !ok || len(when[0]) < 1 {
		DoShutdown(writer, "now")
	} else {
		DoShutdown(writer, when[0])
	}
}

func CancelShutdown(writer http.ResponseWriter, req *http.Request) {
	log.Printf("Got cancel shutdown request")
	DoCancelShutdown(writer, "now")
}

func DoShutdown(writer http.ResponseWriter, when string) {
	log.Printf("Shutting down (%s)\n", when)
	fmt.Fprintf(writer, "Shutting down (%s)\n", when)
	cmd := exec.Command("shutdown", "-h", when)
	out, err := cmd.CombinedOutput()
	if err == nil {
		log.Printf("Shutdown finished successfully:\n%s\n", out)
		fmt.Fprintf(writer, "Shutdown performed successfully:\n%s\n", out)
	} else {
		log.Printf("Error performing shutdown: %v\n%s\n", err, out)
		fmt.Fprintf(writer, "Error performing shutdown: %v\n%s\n", err, out)
	}
}

func DoCancelShutdown(writer http.ResponseWriter, when string) {
	log.Printf("Canceling shutdown\n")
	fmt.Fprintf(writer, "Canceling shutdown\n")
	cmd := exec.Command("shutdown", "-c")
	out, err := cmd.CombinedOutput()
	if err == nil {
		log.Printf("Cancel shutdown finished successfully:\n%s\n", out)
		fmt.Fprintf(writer, "Cancel shutdown performed successfully:\n%s\n", out)
	} else {
		log.Printf("Error canceling shutdown: %v\n%s\n", err, out)
		fmt.Fprintf(writer, "Error canceling shutdown: %v\n%s\n", err, out)
	}
}

func SetupSignals() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTRAP, syscall.SIGKILL)
	go func() {
		s := <-signals
		log.Printf("RECEIVED SIGNAL: %s",s)
		os.Exit(1)
	}()
}

func main() {
	log.Printf("Starting server")
	SetupSignals()

	http.HandleFunc("/shutdown", Shutdown)
	http.HandleFunc("/cancel", CancelShutdown)
	http.ListenAndServe(":8090", nil)
	log.Printf("Stopping server")
}
