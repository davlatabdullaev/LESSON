package main

import (
	"fmt"
	"math/rand"
)

func main() {

// TASK 1

//    num := Task1()
//    fmt.Println(<-num)

// TASK 2

ch := Task2()

for val := range ch {
	fmt.Print(val, ", ")
}

fmt.Println()

}

func Task1() chan int {

ch1 := make(chan int)

num := 0

fmt.Println("Enter number: ")

fmt.Scan(&num)

go func() {

randomNumber := rand.Intn(num)

defer close(ch1)

ch1 <- randomNumber

}()

return ch1

}

func Task2() chan int  {

	ch2 := make(chan int)

	num := 0

	fmt.Println("Enter number: ")

	fmt.Scan(&num)

	go func() {

		defer close(ch2)

		for i:=0; i<num; i++ {

			randNumber := rand.Intn(num)

             ch2 <- randNumber

		}


	}()

return ch2

}