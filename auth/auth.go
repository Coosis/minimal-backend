package auth

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	UserName  string   `bson:"username,omitempty"`
	PswdHash  string   `bson:"pswdhash,omitempty"`
	UserGroup []string `bson:"usergroup,omitempty"`
}

type UserGroup struct {
	GroupName   string   `bson:"groupname,omitempty"`
	Users       []string `bson:"users,omitempty"`
	Permissions []string `bson:"permissions,omitempty"`
}

func Add_user(w http.ResponseWriter, r *http.Request, ctx context.Context, client *mongo.Client) {
	if r.Method != "POST" {
		msg := fmt.Sprintf("Method not allowed: %s, use POST instead", r.Method)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}

	coll := client.Database("db").Collection("users")
	name := r.FormValue("username")
	hash := r.FormValue("pswdhash")
	groups := r.Form["usergroup"]
	user := User{
		UserName: name,
		PswdHash: hash,
		UserGroup: groups,
	}
	_, err := coll.InsertOne(ctx, user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request, ctx context.Context, client *mongo.Client) {
	if r.Method != "POST" {
		msg := fmt.Sprintf("Method not allowed: %s, use POST instead", r.Method)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}

	coll := client.Database("db").Collection("users")
	name := r.FormValue("username")
	hash := r.FormValue("pswdhash")
	var user User
	err := coll.FindOne(ctx, User{UserName: name, PswdHash: hash}).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := gen_token(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(token))
}

func Handle(ctx context.Context, client *mongo.Client) error {
	http.HandleFunc("/auth/add", func(w http.ResponseWriter, r *http.Request) {
		Add_user(w, r, ctx, client)
	})

	http.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		Login(w, r, ctx, client)
	})

	select {
	case <-ctx.Done():
		return ctx.Err()
	}
}
