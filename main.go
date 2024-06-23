package main

import (
	"context"
	"fmt"

	"github.com/Coosis/minimal-backend/auth"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println("Error while connecting to MongoDB!")
		panic(err)
	}
	defer client.Disconnect(context.TODO())
	// auth.CreateAdmin(context.TODO(), client, "__", "1234")
	auth.Handle(context.TODO(), client)
}
