package auth

import (
	"fmt"
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// get permission from a group
func (g *UserGroup) GetPerm(ctx context.Context, client *mongo.Client) ([]string, error) {
	coll := client.Database("db").Collection("groups")
	ug := UserGroup{}
	err := coll.FindOne(ctx, bson.M{"groupname": g.GroupName}).Decode(&ug)
	if err != nil {
		return nil, err
	}

	return ug.Permissions, nil
}

// get permission from a user
func (u *User) GetPerm(ctx context.Context, client *mongo.Client) ([]string, error) {
	coll := client.Database("db").Collection("groups")
	perms := []string{}
	for _, group := range u.UserGroup {
		ug := UserGroup{}
		err := coll.FindOne(ctx, bson.M{"groupname": group}).Decode(&ug)
		if err != nil {
			return nil, err
		}
		perms = append(perms, ug.Permissions...)
	}
	return perms, nil
}

func RightWall(r *http.Request, ctx context.Context, client *mongo.Client, right string) (bool, error) {
	token := r.Header.Get("Authorization")
	token = token[7:]
	if token == "" {
		return false, nil
	}

	// valid token
	username, err := ValidateToken(token)
	if err != nil {
		return false, err
	}

	// get user
	user := User{}
	coll := client.Database("db").Collection("users")
	err = coll.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return false, err
	}

	// admin has all rights
	for _, groupname := range user.UserGroup {
		if groupname == "admin" {
			return true, nil
		}
	}

	// get permissions
	perms, err := user.GetPerm(ctx, client)
	if err != nil {
		return false, err
	}

	// check if user has the right
	for _, perm := range perms {
		if perm == right {
			return true, nil
		}
	}

	return false, nil
}
