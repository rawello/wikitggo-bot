package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"strings"
)

var host = os.Getenv("HOST")
var port = os.Getenv("PORT")
var user = os.Getenv("USER")
var password = os.Getenv("PASSWORD")
var dbname = os.Getenv("DBNAME")
var sslmode = os.Getenv("SSLMODE")

var dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

//Collecting data from bot
func collectData(username string, chatid int64, message string, answer []string) error {

	//Connecting to database
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	//Converting slice with answer to string
	answ := strings.Join(answer, ", ")

	//Creating SQL command
	data := `INSERT INTO users(username, chat_id, message, answer) VALUES($1, $2, $3, $4);`

	//Execute SQL command in database
	if _, err = db.Exec(data, `@`+username, chatid, message, answ); err != nil {
		return err
	}

	return nil
}

//Creating users table in database
func createTable() error {

	//Connecting to database
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	//Creating users Table
	if _, err = db.Exec(`CREATE TABLE users(ID SERIAL PRIMARY KEY, TIMESTAMP TIMESTAMP DEFAULT CURRENT_TIMESTAMP, USERNAME TEXT, CHAT_ID INT, MESSAGE TEXT, ANSWER TEXT);`); err != nil {
		return err
	}

	return nil
}

//Getting number of users who using bot
func getNumberOfUsers() (int64, error) {

	var count int64

	//Connecting to database
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	//Counting number of users
	row := db.QueryRow("SELECT COUNT(DISTINCT username) FROM users;")
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}