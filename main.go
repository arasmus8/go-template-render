package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"src/internal/server"
	"strconv"
)

// InitLogFile initializes the log file
func InitLogFile(logFilePath string) (*os.File, error) {
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("error opening log file: %v", err)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	return logFile, nil
}

func main() {

	logFile, logerr := InitLogFile("./application.log")
	if logerr != nil {
		fmt.Printf("Failed to initialize log file: %v\n", logerr)
		os.Exit(1)
	}
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			log.Println("Error closing Log: ", err)
		}
	}(logFile)

	s := flag.String("port", "8080", "a port in string format") // Default to port 8080 if no arg is passed in

	flag.Parse()

	v, numErr := strconv.Atoi(*s)
	if numErr != nil {
		log.Print("Not a port. Exiting\n")
		return
	}
	log.Printf("Running on %d\n", v)
	newServer := server.NewServer(v, logFile)
	err := newServer.ListenAndServe()
	if err != nil {
		log.Printf("cannot start newServer: %s", err)
	}
}
