package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
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
	count := len(problems)
	if err != nil {
		return err
	}

	var correct int
	r := bufio.NewReader(os.Stdin)

	t := time.NewTimer(time.Duration(limit) * time.Second)
	defer t.Stop()

	go timer(t, &correct, count)

	for _, p := range problems {
		fmt.Printf("%v=", p.question)
		answer, _ := r.ReadString('\n')
		if strings.TrimRight(answer, "\n") == p.answer {
			correct++
		}
	}

	result(correct, count)

	return nil
}

func timer(t *time.Timer, correct *int, count int) {
	<-t.C
	fmt.Println("\nTime is up")
	result(*correct, count)
	os.Exit(0)
}

func result(correct, count int) {
	fmt.Printf("You scored %v out of %v\n", correct, count)
}

//Reads the file and returns an array of problem
func readFile(csvFileName string) ([]problem, error) {

	var problems []problem

	//Open CSV File
	file, err := os.Open(csvFileName)
	if err != nil {
		return problems, err
	}
	//Close file on return
	defer file.Close()

	//Iterate file
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
