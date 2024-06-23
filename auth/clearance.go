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

// given a token, verify that the user is an admin by checking the database
func VerifyAdmin(ctx context.Context, client *mongo.Client, token string) (bool, error) {
	// valid token
	username, err := ValidateToken(token)
	if err != nil {
		return false, err
	}

	// existing user
	exists, user := UserExists(ctx, client, username)
	if !exists {
		return false, nil
	}

	// has admin group
	for _, group := range user.UserGroup {
		if group == "admin" {
			return true, nil
		}
	}
	return false, nil
}

func OnlyAdmin(w http.ResponseWriter, r *http.Request, ctx context.Context, client *mongo.Client) (bool, error) {
	token := r.Header.Get("Authorization")
	fmt.Println(fmt.Sprintf("Token: %s", token))
	if token == "" {
		return false, nil
	}
	// usually the token is in the form "Bearer:<token>"
	// so we remove the first 7 characters
	token = token[7:]

	isAdmin, err := VerifyAdmin(ctx, client, token)
	if err != nil && !isAdmin {
		return false, err
	}

	return true, nil
}

