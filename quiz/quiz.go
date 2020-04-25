package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	c := make(chan int)

	file := flag.String("file", "default", "usage -file filename.csv")

	flag.Parse()
	content := openFile(*file)
	records := csvParser(content)

	r := bufio.NewReader(os.Stdin)
	fmt.Println("\nPress a key to start the quizz !!!!")
	r.ReadString('\n') // not work with all keys

	go askUser(records, c)

	// debug why its blocking todo : unblocking channel

	var score int
	out := false

	// Part 2
	// customize time value with a flag todo
	ch := time.After(10 * time.Second)
	for {
		select {
		case score = <-c:
		case <-ch:
			out = true
			break
		}
		if out {
			break
		}
	}

	result(score, len(records))
	fmt.Println("the end ")

}

func openFile(filename string) string {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(file)
}

func csvParser(content string) [][]string {
	r := csv.NewReader(strings.NewReader(content))
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return records
}

func askUser(records [][]string, c chan int) int {
	goodAnswer := 0
	numberQuestions := len(records)
	r := bufio.NewReader(os.Stdin)

	// ask questions to the user
	for i := 0; i < numberQuestions; i++ {

		fmt.Println(records[i][0] + " ?")
		answer, _ := r.ReadString('\n')
		if strings.TrimRight(answer, "\n") == records[i][1] {
			goodAnswer++

		}
		select {
		case c <- goodAnswer:
		default:
		}

	}

	return goodAnswer

}

func result(score int, numberQuestions int) {
	s1 := fmt.Sprint("Your score is ", score)
	s2 := fmt.Sprint("The quiz has ", numberQuestions, " questions")
	fmt.Println(s1)
	fmt.Println(s2)
}
