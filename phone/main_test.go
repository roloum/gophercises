package main

import (
	"fmt"
	"testing"
)

func TestNormalize(t *testing.T) {

	cases := []struct {
		input    string
		expected string
	}{
		{"1234567890", "1234567890"},
		{"123 456 7891", "1234567891"},
		{"(123) 456 7892", "1234567892"},
		{"(123) 456-7893", "1234567893"},
		{"123-456-7894", "1234567894"},
		{"1234567892", "1234567892"},
		{"(123)456-7892", "1234567892"},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			output := normalize(c.input)
			if output != c.expected {
				t.Errorf("Received: %s, expected: %s", output, c.expected)
			}
		})
	}
	_ = normalize("")
}

func TestConnect(t *testing.T) {
	db, err := connect()
	if err != nil {
		t.Error(err)
	}

	if db.Close() != nil {
		t.Error(err)
	}
}

func TestCRUD(t *testing.T) {
	db, err := connect()
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Connection established")

	defer func() {
		db.Close()
		fmt.Println("Closed connection")
	}()

	tx, err := db.Begin()
	if err != nil {
		t.Error(err)
	}

	defer func() {
		tx.Rollback()
		fmt.Println("Rolled back")
	}()

	number := "9999999999"

	phoneID, err := insert(tx, number)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("ID", phoneID)

	if err := delete(tx, phoneID); err != nil {
		return
	}

	if err := tx.Commit(); err != nil {
		t.Error(err)
	}
	fmt.Println("Committed")

}
