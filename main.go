package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file that contains all the questions to the quiz")
	timeLimit := flag.Int("limit", 30, "time limit for each question")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("File %s failed to be opened!", *csvFilename))

	}
	r := csv.NewReader(file)

	lines, err := r.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Failed to read file!"))
	}
	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.Q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nYou got %d out of %d correct!\n", correct, len(problems))
			return
		case answer := <-answerCh:
			if answer == p.A {
				fmt.Println("Correct!")
				correct++
			} else if answer != p.A {
				fmt.Println("Incorrect!")
			}

		}
	}
	fmt.Printf("You got %d out of %d correct!\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			Q: line[0],
			A: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

type problem struct {
	Q string
	A string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
