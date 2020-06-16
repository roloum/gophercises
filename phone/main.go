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

type phone struct {
	phoneID int64
	number  string
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	db, err := connect()
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}

	phones, err := loadPhones(db)
	if err != nil {
		return err
	}

	for _, p := range phones {

		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer func() {
			tx.Rollback()
		}()

		newNumber := normalize(p.number)
		_, err = getPhoneID(tx, newNumber)
		//Normalized phone does not exist, update row
		if err == sql.ErrNoRows {
			if err = update(tx, p.phoneID, newNumber); err != nil {
				return err
			}
		} else {
			//Normalized phone already exists, delete record
			if err = delete(tx, p.phoneID); err != nil {
				return err
			}
		}
		tx.Commit()
	}

	return nil
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

	result, err := tx.Exec("INSERT INTO "+table+" (number) VALUES (?)", number)
	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

func getPhoneID(tx *sql.Tx, number string) (int64, error) {
	var phoneID int64

	err := tx.QueryRow("SELECT phone_id FROM "+table+" WHERE number=?",
		number).Scan(&phoneID)
	if err != nil {
		return -1, err
	}

	return phoneID, nil
}

func delete(tx *sql.Tx, phoneID int64) error {
	stm, err := tx.Prepare("DELETE FROM " + table + " WHERE phone_id=?")
	if err != nil {
		return err
	}

	_, err = stm.Exec(phoneID)
	return err
}

func update(tx *sql.Tx, phoneID int64, number string) error {
	_, err := tx.Exec("UPDATE phone SET number=? WHERE phone_id=?", number, phoneID)
	return err
}

func loadPhones(db *sql.DB) ([]phone, error) {
	phones := []phone{}

	rows, err := db.Query("SELECT phone_id, number FROM " + table + " ORDER BY 1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			phoneID int64
			number  string
		)
		if err := rows.Scan(&phoneID, &number); err != nil {
			return nil, err
		}

		phones = append(phones, phone{phoneID, number})

	}

	return phones, nil
}
