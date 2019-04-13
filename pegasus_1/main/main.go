package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"pegasus_1/handler"
	"pegasus_1/mySQL"
	"strconv"
)

var mySigningKey = []byte("a")


type Order struct{
	RobotId uint16 `json:"robot_id"`
	Username string `json:"username"`
	Size string `json:"size"`
	ArrivalTime string `json:"arrival"`
	Weight float64 `json:"weight"`
	PickupLoc string `json:"pickup"`
	DropoffLoc string `json:"dropoff"`

}
type resp struct{
	OrderId uint16
}
type arrival struct{
	ArrivalTime string
}

func main() {
	mySQL.NewTable()
	//db, _ := mySQL.Connect()
	//defer db.Close()

	fmt.Println("fire-up-engine")
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(mySigningKey), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	r := mux.NewRouter()
	r.Handle("/signup", http.HandlerFunc(handler.Signup)).Methods("POST", "OPTIONS")
	r.Handle("/login", http.HandlerFunc(handler.Login)).Methods("POST", "OPTIONS")
	// test handle
	r.Handle("/search", jwtMiddleware.Handler(http.HandlerFunc(handler.SearchPath))).Methods("GET", "OPTIONS")
	r.Handle("/testget", jwtMiddleware.Handler(http.HandlerFunc(handler.Test))).Methods("GET", "OPTIONS")
	r.Handle("/order",jwtMiddleware.Handler(http.HandlerFunc(handlerOrder))).Methods("POST","OPTIONS")
	r.Handle("/track",jwtMiddleware.Handler(http.HandlerFunc(handlerTrack))).Methods("GET","OPTIONS")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handlerOrder(w http.ResponseWriter,r *http.Request){
	fmt.Println("Receive order")
	//w.Header().Set("Content-Type","application/json")
	//w.Header().Set("Access-Control-Allow-Origin","*")
	//w.Header().Set("Access-Control-Allow-Headers","Content-Type,Authorization")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "OPTIONS" {
		return
	}
	user:=r.Context().Value("user")
	claims:=user.(*jwt.Token).Claims
	username:=claims.(jwt.MapClaims)["username"]
	fmt.Println(username)
	decoder:=json.NewDecoder(r.Body)
	var p Order
	if err:=decoder.Decode(&p);err!=nil{
		panic(err)
	}
	//fmt.Fprintf(w,"Order received: %d,%s,%f",p.UserId,p.Size,p.Weight)
	db, err := sql.Open("mysql", mySQL.USERNAME + ":" + mySQL.PASSWORD + "@tcp(" +
		mySQL.HOSTNAME + ":" + mySQL.PORT_NUMBER + ")/pegasus")
	fmt.Println("Receive order2")
	if err != nil {
		panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
	}
	// close the db
	defer db.Close()
	querySize:="select count(*) from orders"
	qsz,err :=db.Query(querySize)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Receive order3")
	var sz uint16
	if qsz.Next(){
		if err := qsz.Scan(&sz); err != nil {
			fmt.Println("err", err)
		}
	}
	fmt.Println("sz:" ,sz)
	q,err:=db.Prepare("insert into orders values(?,?,?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
	}
	defer q.Close()
	fmt.Println(p.ArrivalTime)
	fmt.Printf("username: %s, %T\n", username, username)
	_,err=q.Exec(sz + 1, username, p.RobotId, p.Size,p.Weight, p.PickupLoc, p.DropoffLoc, p.ArrivalTime)

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	res:=resp{
		OrderId: sz+1,
	}
	b,err:=json.Marshal(res)
	if err!=nil{
		fmt.Println("error:",err)
	}
	w.Write(b)
}

func handlerTrack(w http.ResponseWriter,r *http.Request){
	fmt.Println("Received one request for track")
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type,Authorization")

	orderId,_:= strconv.ParseInt(r.URL.Query().Get("orderid"),10,32)

	db, err := sql.Open("mysql", mySQL.USERNAME + ":" + mySQL.PASSWORD + "@tcp(" +
		mySQL.HOSTNAME + ":" + mySQL.PORT_NUMBER + ")/pegasus")

	if err != nil {
		panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
	}
	// close the db
	defer db.Close()
	rows,err :=db.Query("select ArrivalTime from orders where OrderId=?", orderId)
	if err != nil {
		fmt.Println("cannot ")
		panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
	}
	defer rows.Close()

	var atime string
	if rows.Next(){

		if err := rows.Scan(&atime); err != nil {
			fmt.Println("err", err)
		}
	}
	res:=arrival{
		ArrivalTime:atime,
	}
	fmt.Println("arrival time:", atime)
	b,err:=json.Marshal(res)
	if err!=nil{
		fmt.Println("error:",err)
	}
	w.Write(b)
}