package main

import "fmt"

func main() {
	ch1 := writeWords()

	newChan := removeRepeated(ch1)

	for word := range newChan {
		fmt.Print(word, ", ")
	}
	fmt.Println()

}

func removeRepeated(ch1 chan string) chan string {
	newChan := make(chan string)

	go func() {
		defer close(newChan)
		previous := ""

		for word := range ch1 {
			if previous != word {
				newChan <- word
			}
			previous = word
		}
	}()
	return newChan
}

func writeWords() chan string {
	words := []string{"apple", "banana", "banana", "banana", "apple", "cherry"}

	ch1 := make(chan string)

	go func() {
		defer close(ch1)
		for _, word := range words {
			ch1 <- word
		}
	}()

	return ch1
}
