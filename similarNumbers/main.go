package main

import (
	"fmt"
	"log"
)

func main() {
	var number int

	fmt.Println("Enter the length of array:")
	fmt.Scan(&number)

	const maxLength = 100

	if number <= 0 || number > maxLength {
		log.Println("wrong array length", maxLength)
		return
	}

	var numbers [maxLength]int

	for i := 0; i < number; i++ {
		fmt.Printf("Enter %d element: ", i+1)
		fmt.Scan(&numbers[i])
	}

	counter := 0
	for i := 0; i < number; i++ {
		for j := i + 1; j < number; j++ {
			if numbers[i] == numbers[j] {
				counter++
			}
		}
	}

	fmt.Println("Teng elementlar soni:", counter)
}
