package auth
import (
	"context"
	"fmt"
	"net/http"
	"encoding/json"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
)
// UserExists: Check whether a user exists
// CreateAdmin: Create an admin user
// AddUser: Add a new user to the database
// DeleteUser: Delete a user from the database
// AddUserToGroup: Add a user to a group
// RemoveUserFromGroup: Remove a user from a group

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
		fmt.Println(err)
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
			fmt.Println("Admin group not found, creating...")
			coll.InsertOne(ctx, UserGroup{GroupName: "admin"})
		} else {
			fmt.Println(err)
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
		fmt.Println(err)
		return err
	}

	// add user to admin group
	coll = client.Database("db").Collection("groups")
	query := bson.M{"$addToSet": bson.M{"users": username}}
	_, err2 := coll.UpdateOne(ctx, bson.M{"groupname": "admin"}, query)
	if err2 != nil {
		fmt.Println(err2)
		return err2
	}

	fmt.Println("Admin user created: " + username)
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

	//remove in production
	fmt.Println("User added:")
	fmt.Println(user)

	for _, group := range groups {
		AddUserToGroup(w, r, ctx, client, group, name)
		//remove in production
		fmt.Println("User added to group: " + group)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := Response{Message: "User added"}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func DeleteUser(w http.ResponseWriter, r *http.Request, ctx context.Context, client *mongo.Client, username string) {
	if r.Method != "DELETE" {
		msg := fmt.Sprintf("Method not allowed: %s, use DELETE instead", r.Method)
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

	//remove in production
	fmt.Println("User deleted: " + username)

	response := Response{Message: "User deleted"}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
