package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserExists: Check whether a user exists
// CreateAdmin: Create an admin user
// AddUser: Add a new user to the database
// DeleteUser: Delete a user from the database

type Response struct {
	Message string `json:"message"`
}

// whether a user exists, if so return the user
func UserExists(ctx context.Context, client *mongo.Client, username string) (bool, User) {
	coll := client.Database("db").Collection("users")
	var user User
	err := coll.FindOne(ctx, User{UserName: username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, user
		}

		log.Println(err)
	}
	return true, user
}

// create an admin user
func CreateAdmin(ctx context.Context, client *mongo.Client, w http.ResponseWriter, r *http.Request) error {
	username := r.FormValue("username")
	hash := r.FormValue("pswdhash")

	coll := client.Database("db").Collection("groups")

	// create admin group if it doesn't exist
	singleResult := coll.FindOne(ctx, bson.M{"groupname": "admin"})
	err := singleResult.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("Admin group not found, creating...")
			coll.InsertOne(ctx, UserGroup{GroupName: "admin"})
		} else {
			log.Println(err)
			return err
		}
	}

	//stop if user already exists
	exists, _ := UserExists(ctx, client, username)
	if exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return nil
	}

	// create admin user
	user := User{
		UserName:  username,
		PswdHash:  hash,
		UserGroup: []string{"admin"},
	}
	coll = client.Database("db").Collection("users")
	_, err = coll.InsertOne(ctx, user)
	if err != nil {
		log.Println(err)
		return err
	}

	// add user to admin group
	coll = client.Database("db").Collection("groups")
	query := bson.M{"$addToSet": bson.M{"users": username}}
	_, err2 := coll.UpdateOne(ctx, bson.M{"groupname": "admin"}, query)
	if err2 != nil {
		log.Println(err2)
		return err2
	}

	log.Printf("Admin user %s created\n", username)
	return nil
}

func AddUser(w http.ResponseWriter, r *http.Request, ctx context.Context, client *mongo.Client) {
	if r.Method != "POST" {
		msg := fmt.Sprintf("Method not allowed: %s, use POST instead", r.Method)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}

	// check whether user already exists
	coll := client.Database("db").Collection("users")
	name := r.FormValue("username")
	if name == "" {
		http.Error(w, "Username not provided", http.StatusBadRequest)
		return
	}

	exists, _ := UserExists(ctx, client, name)
	if exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	hash := r.FormValue("pswdhash")
	if hash == "" {
		http.Error(w, "Password hash not provided", http.StatusBadRequest)
		return
	}
	groups := r.Form["usergroup"]
	user := User{
		UserName:  name,
		PswdHash:  hash,
		UserGroup: groups,
	}
	_, err := coll.InsertOne(ctx, user)

	log.Printf("User %s added\n", name)

	for _, group := range groups {
		AddUserToGroup(w, r, ctx, client, group, name)
		log.Printf("User %s added to group %s\n", name, group)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := Response{Message: fmt.Sprintf("User %s added", name)}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func DeleteUser(w http.ResponseWriter, r *http.Request, ctx context.Context, client *mongo.Client, username string) {
	if r.Method != "POST" {
		msg := fmt.Sprintf("Method not allowed: %s, use POST instead", r.Method)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}

	coll := client.Database("db").Collection("users")
	_, err := coll.DeleteOne(ctx, bson.M{"username": username})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	coll = client.Database("db").Collection("groups")
	removal := bson.M{"$pull": bson.M{"users": username}}
	_, err = coll.UpdateMany(ctx, bson.M{"users": username}, removal)

	log.Printf("User %s deleted\n", username)

	response := Response{Message: fmt.Sprintf("User %s deleted", username)}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
