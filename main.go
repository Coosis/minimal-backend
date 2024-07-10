package main

import (
	"context"

	"os"
	"os/signal"
	"syscall"

	"github.com/Coosis/minimal-backend/auth"
	l "github.com/Coosis/minimal-backend/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	go l.LogHandler()
	l.Logchan <- "Starting server..."

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		l.Logchan <- "Error while connecting to MongoDB!"
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		auth.Handle(context.TODO(), client)
	}()

	<-sigc
	l.Logchan <- "Shutting down server..."
}
