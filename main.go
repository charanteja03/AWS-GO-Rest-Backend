package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	//envs
	_ "sfr-backend/docs"
	"sfr-backend/server"

	"github.com/joho/godotenv"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	log "github.com/sirupsen/logrus"
)

func setupGracefulShutdown() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Info(" got interrupt signal: ", <-sigChan)
		log.Info("Shutting down....")
		os.Exit(0)
	}()
}

func initLog() {
	maxAge := 0
	maxAgeString := os.Getenv("LOG_MAX_AGE")
	if maxAgeString == "" {
		maxAge = 60
	} else {
		maxAge, _ = strconv.Atoi(maxAgeString)
	}

	rotationTime := 0
	rotationTimeString := os.Getenv("LOG_ROTATION_TIME")
	if rotationTimeString == "" {
		rotationTime = 10
	} else {
		rotationTime, _ = strconv.Atoi(rotationTimeString)
	}

	logPath := os.Getenv("LOGPATH")
	if logPath == "" {
		logPath = "logs/"
	}

	path := os.Getenv("LOGPATH") + "sfservice_log"
	writer, err := rotatelogs.New(
		fmt.Sprintf("%s.%s", path, "%Y_%m_%d_%H_%M_%S"),
		rotatelogs.WithLinkName(logPath+"current"),
		rotatelogs.WithMaxAge(time.Second*time.Duration(maxAge)),
		rotatelogs.WithRotationTime(time.Second*time.Duration(rotationTime)),
	)
	if err != nil {
		log.Fatalf("Failed to Initialize Log File %s", err)
	}
	log.SetOutput(writer)
}

func main() {
	// Load the .env file in the current directory
	godotenv.Load()
	initLog()
	setupGracefulShutdown()
	// Start sever
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8181"
	}
	http.ListenAndServe(":"+serverPort, server.StartServer())
}
