package handler

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"pegasus_1/mySQL"
	"pegasus_1/utils"
	"time"
)

var mySigningKey = []byte("a")

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one login request")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		return
	}

	// decode request, check err
	decoder := json.NewDecoder(r.Body)
	var user utils.User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "cannot decode user data from client", http.StatusBadRequest)
		fmt.Println("cannot decode user data from client")
		return
	}

	// connect to database
	db, err := mySQL.Connect()
	defer db.Close()
	// check connection failure, incorrect pass or userid errors
	if err != nil {
		http.Error(w, "DB Connection failed", http.StatusInternalServerError)
		fmt.Println("DB Connection failed")
	}
	if err = db.CheckUser(user); err != nil {
		if err.Error() == "incorrect password or user_id" {
			http.Error(w, "incorrect password or user_id", http.StatusUnauthorized)
			fmt.Println("wrong password or user_id")

		} else {
			http.Error(w, "login fail", http.StatusInternalServerError)
			fmt.Println("whuat?")
		}
		return
	}

	// if login, create token and check building token error
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id" : user.UserID,
		"username" : user.Username,
		"exp" : time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		http.Error(w, "cannot generate token", http.StatusInternalServerError)
		fmt.Println("cannot generate token")
		return
	}

	// send token back
	w.Write([]byte(tokenString))
}