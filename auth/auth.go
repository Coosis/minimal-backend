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

func UserExists(ctx context.Context, client *mongo.Client, username string) bool {
	coll := client.Database("db").Collection("users")
	var user User
	err := coll.FindOne(ctx, User{UserName: username}).Decode(&user)
	return err == nil
}

func AddUserToGroup(w http.ResponseWriter, r *http.Request, ctx context.Context, client *mongo.Client, groupname string, username string) {
	if r.Method != "POST" {
		msg := fmt.Sprintf("Method not allowed: %s, use POST instead", r.Method)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}

	coll := client.Database("db").Collection("usergroups")
	var ug UserGroup
	coll.FindOne(ctx, UserGroup{GroupName: groupname}).Decode(&ug)
	_, err := coll.UpdateOne(ctx, UserGroup{GroupName: groupname}, UserGroup{Users: append(ug.Users, username)})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AddUser(w http.ResponseWriter, r *http.Request, ctx context.Context, client *mongo.Client) {
	if r.Method != "POST" {
		msg := fmt.Sprintf("Method not allowed: %s, use POST instead", r.Method)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}

	coll := client.Database("db").Collection("users")
	name := r.FormValue("username")
	if UserExists(ctx, client, name) {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	hash := r.FormValue("pswdhash")
	groups := r.Form["usergroup"]
	user := User{
		UserName: name,
		PswdHash: hash,
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

	token, err := Gen_token(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(token))
	fmt.Println("User logged in: " + name)
}

func Handle(ctx context.Context, client *mongo.Client) error {
	http.HandleFunc("/auth/add", func(w http.ResponseWriter, r *http.Request) {
		AddUser(w, r, ctx, client)
	})

	http.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		Login(w, r, ctx, client)
	})

	http.HandleFunc("/auth/addtogroup", func(w http.ResponseWriter, r *http.Request) {
		AddUserToGroup(w, r, ctx, client, r.FormValue("groupname"), r.FormValue("username"))
	})

	http.ListenAndServe(":8080", nil)

	select {
	case <-ctx.Done():
		return ctx.Err()
	}
}
