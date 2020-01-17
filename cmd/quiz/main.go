package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

type problem struct {
	question string
	answer   string
}

func init() {
	rand.Seed(time.Now().UnixNano())
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
	var shuffle bool

	flag.StringVar(&csvFileName, "csv", "problems.csv",
		"A CSV in the format of 'question,answer'")
	flag.IntVar(&limit, "limit", 30, "The time limit for the quiz in seconds")
	flag.BoolVar(&shuffle, "s", false, "Shuffle problems")
	flag.Parse()

	//Read file in problem array
	problems, err := readFile(csvFileName)
	if err != nil {
		return err
	}

	if shuffle {
		for i := range problems {
			j := rand.Intn(i + 1)
			problems[i], problems[j] = problems[j], problems[i]
		}
	}

	var correct int
	t := time.NewTimer(time.Duration(limit) * time.Second)

problemLoop:
	for _, p := range problems {
		fmt.Printf("%v=", p.question)

		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- strings.TrimSpace(answer)
		}()

		select {
		case answer := <-answerCh:
			if answer == p.answer {
				correct++
			}
		case <-t.C:
			fmt.Println("")
			break problemLoop
		}

	}

	fmt.Printf("You scored %v out of %v\n", correct, len(problems))

	return nil
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
