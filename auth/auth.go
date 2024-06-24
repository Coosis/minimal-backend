package auth

import (
	"context"
	"fmt"
	"net/http"

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
			fmt.Println(err)
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
			fmt.Println(err)
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
			fmt.Println(err)
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
			fmt.Println(err)
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
			fmt.Println(err)
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
