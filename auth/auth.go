package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

// Login: Log in a user and return a jwt token
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

type LoginResponse struct {
	UserName string `json:"username"`
	Token    string `json:"token"`
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

	response := LoginResponse{UserName: user.UserName, Token: token}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	log.Printf("User %s logged in\n", user.UserName)
}

func Handle(ctx context.Context, client *mongo.Client) error {
	go func() {
		http.HandleFunc("/auth/createadmin", func(w http.ResponseWriter, r *http.Request) {
			err := CreateAdmin(ctx, client, w, r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
		http.ListenAndServe("http://localhost:8081", nil)
	}()

	http.HandleFunc("/auth/add", func(w http.ResponseWriter, r *http.Request) {
		AddUser(w, r, ctx, client)
	})

	http.HandleFunc("/auth/del", func(w http.ResponseWriter, r *http.Request) {
		cle, err := RightWall(r, ctx, client, "delete_user")
		if !cle {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		DeleteUser(w, r, ctx, client, r.FormValue("username"))
	})

	http.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		Login(w, r, ctx, client)
	})

	http.HandleFunc("/auth/addtogroup", func(w http.ResponseWriter, r *http.Request) {
		cle, err := RightWall(r, ctx, client, "edit_group")
		if !cle {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		AddUserToGroup(w, r, ctx, client, r.FormValue("groupname"), r.FormValue("username"))
	})

	http.HandleFunc("/auth/rmfromgroup", func(w http.ResponseWriter, r *http.Request) {
		cle, err := RightWall(r, ctx, client, "edit_group")
		if !cle {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		RemoveUserFromGroup(w, r, ctx, client, r.FormValue("groupname"), r.FormValue("username"))
	})

	http.HandleFunc("/auth/addrighttogroup", func(w http.ResponseWriter, r *http.Request) {
		cle, err := RightWall(r, ctx, client, "edit_group")
		if !cle {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		AddRightToGroup(w, r, ctx, client, r.FormValue("groupname"), r.FormValue("right"))
	})

	http.HandleFunc("/auth/rmrightfromgroup", func(w http.ResponseWriter, r *http.Request) {
		cle, err := RightWall(r, ctx, client, "edit_group")
		if !cle {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		RemoveRightFromGroup(w, r, ctx, client, r.FormValue("groupname"), r.FormValue("right"))
	})

	http.ListenAndServe(":8080", nil)

	select {
	case <-ctx.Done():
		return ctx.Err()
	}
}
