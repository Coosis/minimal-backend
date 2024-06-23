package auth

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateAdmin: Create an admin user
// UserExists: Check whether a user exists
// AddUser: Add a new user to the database
// Login: Log in a user and return a jwt token
// AddUserToGroup: Add a user to a group

// Handle: listen and serve

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

func CreateAdmin(ctx context.Context, client *mongo.Client, username string, hash string) error {
	coll := client.Database("db").Collection("groups")

	// create admin group if it doesn't exist
	ug := UserGroup{}
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

	// stop if admin group already has users
	singleResult.Decode(&ug)
	if ug.Users != nil {
		fmt.Println("Admin already exists")
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
	return nil
}

// whether a user exists, if so return the user
func UserExists(ctx context.Context, client *mongo.Client, username string) (bool, User) {
	coll := client.Database("db").Collection("users")
	var user User
	err := coll.FindOne(ctx, User{UserName: username}).Decode(&user)
	return err == nil, user
}

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
	exists, _ := UserExists(ctx, client, name)
	if exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	hash := r.FormValue("pswdhash")
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

	token, err := GenToken(name)
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
		admin, err := OnlyAdmin(w, r, ctx, client)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !admin {
			fmt.Println(err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		AddUserToGroup(w, r, ctx, client, r.FormValue("groupname"), r.FormValue("username"))
	})

	http.ListenAndServe(":8080", nil)

	select {
	case <-ctx.Done():
		return ctx.Err()
	}
}
