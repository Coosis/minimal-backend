package auth
import (
	"context"
	"fmt"
	"log"
	"net/http"
	"encoding/json"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
)
// AddUserToGroup: Add a user to a group
// RemoveUserFromGroup: Remove a user from a group

func AddUserToGroup(w http.ResponseWriter, r *http.Request, ctx context.Context, client *mongo.Client, groupname string, username string) {
	if r.Method != "POST" {
		msg := fmt.Sprintf("Method not allowed: %s, use POST instead", r.Method)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}

	coll := client.Database("db").Collection("groups")
	// load usergroup
	var ug UserGroup
	err := coll.FindOne(ctx, bson.M{"groupname": groupname}).Decode(&ug)
	// create usergroup if it doesn't exist
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ug = UserGroup{GroupName: groupname, Users: []string{username}, Permissions: []string{}}
			coll.InsertOne(ctx, ug)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// add user to group
	addition := bson.M{"$addToSet": bson.M{"users": username}}
	_, err = coll.UpdateOne(ctx, bson.M{"groupname": groupname}, addition)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// add group to user
	coll = client.Database("db").Collection("users")
	addition = bson.M{"$addToSet": bson.M{"usergroup": groupname}}
	_, err = coll.UpdateOne(ctx, bson.M{"username": username}, addition)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("User %s added to group %s\n", username, groupname)

	response := Response{Message: fmt.Sprintf("User %s added to group %s", username, groupname)}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func RemoveUserFromGroup(w http.ResponseWriter, r *http.Request, ctx context.Context, client *mongo.Client, groupname string, username string) {
	if r.Method != "POST" {
		msg := fmt.Sprintf("Method not allowed: %s, use POST instead", r.Method)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}

	coll := client.Database("db").Collection("groups")
	// remove user from group
	removal := bson.M{"$pull": bson.M{"users": username}}
	_, err := coll.UpdateOne(ctx, bson.M{"groupname": groupname}, removal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// remove group from user
	coll = client.Database("db").Collection("users")
	removal = bson.M{"$pull": bson.M{"usergroup": groupname}}
	_, err = coll.UpdateOne(ctx, bson.M{"username": username}, removal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("User %s removed from group %s\n", username, groupname)

	response := Response{Message: fmt.Sprintf("User %s removed from group %s", username, groupname)}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func AddRightToGroup(w http.ResponseWriter, r *http.Request, ctx context.Context, client *mongo.Client, groupname string, right string) {
	if r.Method != "POST" {
		msg := fmt.Sprintf("Method not allowed: %s, use POST instead", r.Method)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}

	coll := client.Database("db").Collection("groups")
	// add right to group
	addition := bson.M{"$addToSet": bson.M{"permissions": right}}
	_, err := coll.UpdateOne(ctx, bson.M{"groupname": groupname}, addition)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Right %s added to group: %s\n", right, groupname)

	response := Response{Message: fmt.Sprintf("Right %s added to group %s", right, groupname)}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func RemoveRightFromGroup(w http.ResponseWriter, r *http.Request, ctx context.Context, client *mongo.Client, groupname string, right string) {
	if r.Method != "POST" {
		msg := fmt.Sprintf("Method not allowed: %s, use POST instead", r.Method)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}

	coll := client.Database("db").Collection("groups")
	// remove right from group
	removal := bson.M{"$pull": bson.M{"permissions": right}}
	_, err := coll.UpdateOne(ctx, bson.M{"groupname": groupname}, removal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Right %s removed from group: %s\n", right, groupname)

	response := Response{Message: fmt.Sprintf("Right %s removed from group %s", right, groupname)}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
