package main

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
)

const (
	host   = "127.0.0.1"
	port   = "3306"
	user   = "root"
	pass   = "123456"
	dbname = "gophercises"
	table  = "phone"
)

func main() {
	db, err := connect()
	if err != nil {
		er(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	// result, err := db.Query("SELECT * FROM phone")
	// if err != nil {
	// 	fmt.Printf("Error: %s", err)
	// }
	// fmt.Println(result)
}

func normalize(phone string) string {
	re := regexp.MustCompile(`[^0-9]`)
	return re.ReplaceAllString(phone, "")
}

func connect() (*sql.DB, error) {
	return sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		user, pass, host, port, dbname))
}

func insert(tx *sql.Tx, number string) (int64, error) {

	stm, err := tx.Prepare("INSERT INTO " + table + " (number) VALUES (?)")
	if err != nil {
		return -1, err
	}

	result, err := stm.Exec(number)
	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

func er(e error) {
	fmt.Println(e)
	os.Exit(1)
}
