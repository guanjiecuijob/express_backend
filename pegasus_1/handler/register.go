package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pegasus_1/mySQL"
	"pegasus_1/utils"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one signup request")
	w.Header().Set("Content-Type", "text/palin")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "OPTIONS" {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var user utils.User
	// handle err by informal user type
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "cannot decode user data from client", http.StatusBadRequest)
		fmt.Printf("cannot decode user data from client %v.\n", err)
		return
	}
	if user.UserID == "" || user.Password == "" {
		http.Error(w,"Invalid username or password", http.StatusBadRequest)
		fmt.Printf("Invalid username or password.\n")
		return
	}
	db, err := mySQL.Connect()
	defer db.Close()

	// handle err generated by inserting user into database
	if err != nil {
		http.Error(w, "DB Connection failed", http.StatusInternalServerError)
		return
	}
	if err = db.AddUser(user); err != nil {
		if err.Error() == "Users already exists" {
			http.Error(w, "User already exists", http.StatusBadRequest)
		} else {
			http.Error(w, "failed to save to DB", http.StatusInternalServerError)
		}
		return
	}

	w.Write([]byte("User added succussfully."))
}
