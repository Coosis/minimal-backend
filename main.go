package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Coosis/minimal-auth/auth"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("http://localhost:27017"))
	if err != nil {
		fmt.Println("Error while connecting to MongoDB!")
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	http.ListenAndServe(":8080", nil)
}
