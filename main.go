package main

import (
	"context"
	"log"

	"os"
	"os/signal"
	"syscall"

	"github.com/Coosis/minimal-backend/logger"
	"github.com/Coosis/minimal-backend/auth"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	log.Println("Starting server...")
	defer logger.LogFile.Close()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Println("Error while connecting to MongoDB!")
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		auth.Handle(context.TODO(), client)
	}()

	<-sigc
	log.Println("Shutting down server...")
}
