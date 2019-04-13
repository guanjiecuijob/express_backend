package mySQL

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"pegasus_1/utils"
)

type DB struct {
	db *sql.DB
}

func Connect() (DB, error) {
	temp, err := sql.Open("mysql", USERNAME + ":" + PASSWORD + "@tcp(" +
		HOSTNAME + ":" + PORT_NUMBER + ")/pegasus")

	db := DB{temp}

	return db, err
}

func (conn *DB) Close() {
	conn.db.Close()
}

func (conn *DB) AddUser(user utils.User) error {
	q := "insert ignore into users values (?, ?, ?)"
	stmt, err := conn.db.Prepare(q)
	defer stmt.Close()
	res, err := stmt.Exec(user.UserID, user.Password, user.Username)
	if err != nil {
		return err
	}
	// check if insert successful
	if num, _ := res.RowsAffected(); num == 0 {
		return errors.New("Users already exists")
	}
	fmt.Printf("User is add: %s.\n", user.Username)
	return nil
}

func (conn *DB) CheckUser(user utils.User) error {
	q := "select password from users where user_id = ?"
	rows, err := conn.db.Query(q, user.UserID)
	defer rows.Close()
	if err != nil {
		return err
	}

	var password string

	if !rows.Next() {
		return errors.New("incorrect password or user_id")
	}
	rows.Scan(&password)
	fmt.Printf("login: %s.\n", user.Username)
	return nil
}

func (conn *DB) Get(userID string) string {

	q := "select * from users where user_id = ?"

	query, _ := conn.db.Query(q, userID)
	defer query.Close()

	var name string
	var name2 string
	var name3 string

	for query.Next() {
		if err := query.Scan(&name, &name2, &name3); err != nil {
			fmt.Println("err", err)
		}
		fmt.Println(name, name2, name3, "end")
	}
	return fmt.Sprintf("ID: %s.\nName: %s .", name, name3)
}






