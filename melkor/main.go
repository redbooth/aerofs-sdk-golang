package main

// The entrypoint for the Melkor webapp demonstrating the AeroFS Golang SDK

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

// Global logger
var logger *log.Logger

// Server <host>:<port>
var hostName string

var appConfig string

func main() {

	// Parse CLI arguments
	if len(os.Args[1:]) != 3 {
		fmt.Println("Not enough arguments : ./melkor <host> <ip> <appconfig.json>")
		os.Exit(1)
	}
	host := os.Args[1]
	port := os.Args[2]
	appConfig = os.Args[3]
	hostName = fmt.Sprintf("%s:%s", host, port)

	// Initialize logger
	err := initLogger()
	if err != nil {
		fmt.Println("Unable to initialize log file")
		os.Exit(1)
	}
	logger.Print("Melkor beginning startup...")

	// Set Handlers
	router := mux.NewRouter()

	// Static fileserver for css
	resHandler := http.FileServer(http.Dir("./resources/"))
	http.Handle("/resources/", http.StripPrefix("/resources/", resHandler))

	// Authentication
	router.HandleFunc("/tokenization", tokenization).Methods("GET")
	router.HandleFunc("/", defaultHandler).Methods("GET")
	router.HandleFunc("/login", loginEntryHandler).Methods("GET")
	router.HandleFunc("/login", loginSubmitHandler).Methods("POST")

	// View Pages
	router.HandleFunc("/devices", yourDevicesHandler).Methods("GET")
	router.HandleFunc("/files", yourFilesHandler).Methods("GET")
	router.HandleFunc("/totalusers", totalUsersHandler).Methods("GET")
	http.Handle("/", router)

	http.ListenAndServe(hostName, nil)
}

// Initialize the Global server logger
func initLogger() error {
	t := time.Now()
	logTime := fmt.Sprintf("%d-%d-%d_%d-%d-%d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	logName := fmt.Sprintf("Melkor_Logs_%s", logTime)
	logFile, err := os.OpenFile(logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	// Log the file location, time and date
	logger = log.New(logFile, "", log.LstdFlags|log.Lshortfile)

	return nil
}
