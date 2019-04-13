package handler

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"pegasus_1/mySQL"
	"pegasus_1/utils"
)

func Test(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one test request")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		return
	}
	u:=r.Context().Value("user")
	claims:=u.(*jwt.Token).Claims
	username:=claims.(jwt.MapClaims)["username"]

	fmt.Println("this is the username decoded: " , username)
	decoder := json.NewDecoder(r.Body)
	var user utils.User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "cannot decode user data from client", http.StatusBadRequest)
		fmt.Println("cannot decode user data from client")
		return
	}

	db, err := mySQL.Connect()
	defer db.Close()

	if err != nil {
		http.Error(w, "DB Connection failed", http.StatusInternalServerError)
		fmt.Println("DB Connection failed")
	}

	res := db.Get(user.UserID)

	w.Write([]byte(res))
}