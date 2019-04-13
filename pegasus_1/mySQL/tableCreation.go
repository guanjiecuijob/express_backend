package mySQL

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// delete elder pegasus mySQL and create a new one
func NewTable() {
	fmt.Println("start to create a new pegasus datasbase ......")
	db, err := sql.Open("mysql", USERNAME + ":" + PASSWORD + "@tcp(" +
		HOSTNAME + ":" + PORT_NUMBER + ")/") // mysqlroot
	// close the db
	//sql.Open("mysql", "root@cloudsql(project-id:instance-name)/dbname")
	defer db.Close()
	checkErr(err)

	name := "pegasus"
	_, err = db.Exec("CREATE DATABASE " + name)
	checkErr(err)

	_, err = db.Exec("USE " + name)
	checkErr(err)

	q := "create table users (" +
		"user_id varchar(255) not null," +
		"password varchar(255) not null," +
		"username varchar(255) not null," +
		"primary key (user_id)" +
		")"
	_, err = db.Exec(q)
	checkErr(err)

	q = "create table orders (" +
		"OrderId INT," +
		"Username VARCHAR(255)," +
		"RobotId INT," +
		"Size VARCHAR(50)," +
		"Weight DECIMAL," +
		"PickupLocation VARCHAR(255)," +
		"DROPoffLocation VARCHAR(255)," +
		"ArrivalTime VARCHAR(50)," +
		"PRIMARY KEY(OrderId)" +
		")"

	_, err = db.Exec(q)
	checkErr(err)

	fmt.Println("database pegasus created!")
	//_,err = db.Exec("CREATE DATABASE pegasusb00")
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	fmt.Println("Successfully created database..")
	//}
	//// drop pre pegasus database
	//q := "drop database if exists pegasus"
	//db.Query(q)
	//
	//// create users table

}

func checkErr(err error) {
	if (err != nil) {
		panic(err)
	}
}
