package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// problem struct holds a question and its answer
type problem struct {
	question string
	answer   string
}

func main() {
	// parse command line flags
	csvFileName, timeLimit, shuffle := parseFlags()

	// open the CSV file
	file, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV File: %s\n", *csvFileName))
	}
	defer func() {
		if err := file.Close(); err != nil {
			exit(fmt.Sprintf("Failed to close the CSV File: %s\n", *csvFileName))
		}
	}()

	// Read and parse the CSV file
	problems, err := readCSV(file)
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}

	// Shuffle problems if the shuffle flag is set
	if *shuffle {
		shuffleProblems(problems)
	}

	// Run the quiz
	runQuiz(problems, *timeLimit)
}

// parseFlags parses command line flags for CSV file name and time limit
func parseFlags() (*string, *int, *bool) {
	csvFileName := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	shuffle := flag.Bool("shuffle", false, "shuffle the order of the quiz questions")
	flag.Parse()
	return csvFileName, timeLimit, shuffle
}

// readCSV reads the CSV file and returns a slice of problems
func readCSV(file *os.File) ([]problem, error) {
	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	return parseLines(records), nil
}

// parseLines reads the CSV lines and returns a slice of problems
func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))

	for i, line := range lines {
		ret[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}

	return ret
}

// shuffleProblems shuffles the order of the quiz questions
func shuffleProblems(problems []problem) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(problems), func(i, j int) {
		problems[i], problems[j] = problems[j], problems[i]
	})
}

// runQuiz conducts the quiz with the provided problems and time limit
func runQuiz(problems []problem, timeLimit int) {
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
	correct := 0

problemLoop:
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.question)

		answerCh := make(chan string)
		go getAnswer(answerCh)

		select {
		case <-timer.C:
			fmt.Println("Time is up!")
			break problemLoop
		case ans := <-answerCh:
			if ans == p.answer {
				correct++
			}
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
	os.Exit(0)
}

// getAnswer reads an answer from the user and sends it to the provided channel
func getAnswer(answerCh chan string) {
	var ans string
	if _, err := fmt.Scanf("%s\n", &ans); err != nil {
		fmt.Println("Failed to read answer.")
		answerCh <- ""
		return
	}
	answerCh <- ans
}

// exit prints an error message and exits the program
func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
