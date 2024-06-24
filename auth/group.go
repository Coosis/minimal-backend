package auth
import (
	"context"
	"fmt"
	"net/http"

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

	//remove in production
	fmt.Println("User added to group: " + groupname)
}

func RemoveUserFromGroup(w http.ResponseWriter, r *http.Request, ctx context.Context, client *mongo.Client, groupname string, username string) {
	if r.Method != "DELETE" {
		msg := fmt.Sprintf("Method not allowed: %s, use DELETE instead", r.Method)
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

	//remove in production
	fmt.Println("User removed from group: " + groupname)
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

	//remove in production
	fmt.Println("Right added to group: " + groupname)
}

func RemoveRightFromGroup(w http.ResponseWriter, r *http.Request, ctx context.Context, client *mongo.Client, groupname string, right string) {
	if r.Method != "DELETE" {
		msg := fmt.Sprintf("Method not allowed: %s, use DELETE instead", r.Method)
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

	//remove in production
	fmt.Println("Right removed from group: " + groupname)
}
