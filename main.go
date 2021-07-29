package main

import (
	"com.aharakitchen/app/database"
	"com.aharakitchen/app/events"
	"com.aharakitchen/app/router"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
)

func init() {
	// create database connection instance for first time
	go events.KafkaConsumerGroup()
	database.ConnectToRedis()
	_, err := database.ConnectToDB()
	if err != nil {
		return 
	}
}

func main() {
	fmt.Println(runtime.NumGoroutine())
	app := router.Setup()

	// graceful shutdown on signal interrupts
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		_ = <-c
		fmt.Println("Shutting down...")
		fmt.Println(runtime.NumGoroutine())
		_ = app.Shutdown()
	}()

	if err := app.Listen(":8080"); err != nil {
		log.Panic(err)
	}
}