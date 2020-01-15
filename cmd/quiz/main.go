package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
)

type problem struct {
	question string
	answer   string
}

func main() {

	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func run() error {

	var csvFileName string
	var limit int

	flag.StringVar(&csvFileName, "csv", "problems.csv",
		"A CSV in the format of 'question,answer'")
	flag.IntVar(&limit, "limit", 30, "The time limit for the quiz in seconds")
	flag.Parse()

	//Read file in problem array
	problems, err := readFile(csvFileName)
	if err != nil {
		return err
	}
	fmt.Println(problems)

	return nil
}

//Reads the file and returns an array of problem
func readFile(csvFileName string) ([]problem, error) {

	var problems []problem

	file, err := os.Open(csvFileName)
	if err != nil {
		return problems, err
	}
	defer file.Close()

	r := csv.NewReader(file)
	for {
		var row []string
		row, err = r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return problems, err
		}
		problems = append(problems, problem{row[0], row[1]})
	}

	return problems, nil
}
