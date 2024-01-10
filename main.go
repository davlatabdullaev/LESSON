package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

func main() {
	// EXERCISE 1

	// var wg sync.WaitGroup
	// ch := make(chan int)
	// wg.Add(1)
	// go Exercise1(ch, &wg)
	// go func() {
	// 	for num := range ch {
	// 		fmt.Println(num)
	// 	}
	// }()
	// wg.Wait()

	// EXERCISE 2

	// var wg sync.WaitGroup
	// ch := make(chan int)
	// wg.Add(1)
	// go Exercise2(ch, &wg)
	// fmt.Println(<-ch)
	// wg.Wait()

	// EXERCISE 3

	// var wg sync.WaitGroup
	// ch := make(chan int)
	// wg.Add(1)
	// go Exercise3(ch, &wg)
	// fmt.Println(<-ch)
	// wg.Wait()

	// EXERCISE 4

	// var wg sync.WaitGroup
	// ch := make(chan string)
	// wg.Add(1)
	// go Exercise4(ch, &wg)
	// fmt.Println(<-ch)
	// wg.Wait()

	//EXERCISE 5

	// var wg sync.WaitGroup
	// ch := make(chan int)
	// wg.Add(1)
	// go Exercise5(ch, &wg)
	// fmt.Println(<-ch)
	// wg.Wait()

	//EXERCISE 6

	// var wg sync.WaitGroup
	// ch := make(chan string)
	// wg.Add(1)

	// go Exercise6(ch, &wg)

	// fmt.Println(<-ch)
	// wg.Wait()

	//EXERCISE 7

	// var wg sync.WaitGroup
	// ch := make(chan []int)
	// wg.Add(1)

	// go Exercise7(ch, &wg)

	// for val :=range ch {
	// 	fmt.Println(val)
	// }
	// wg.Wait()

	//EXERCISE 8

	// var wg sync.WaitGroup
	// ch := make(chan int)
	// wg.Add(1)

	// go Exercise8(ch, &wg)

	// fmt.Println(<-ch)

	// wg.Wait()

	//EXERCISE 9

	// var wg sync.WaitGroup
	// ch := make(chan []int)
	// wg.Add(1)

	// go Exercise9(ch, &wg)

	// for val :=range ch{
	// 	fmt.Println(val)
	// }

	// wg.Wait()

	//EXERCISE 10

	// var wg sync.WaitGroup
	// ch := make(chan int)
	// wg.Add(1)

	// go Exercise10(ch, &wg)

	// fmt.Println(<-ch)

	// wg.Wait()

	//EXERCISE 11

	// var wg sync.WaitGroup
	// ch := make(chan int)
	// wg.Add(1)

	// go Exercise11(ch, &wg)

	// fmt.Println(<-ch)

	// wg.Wait()

	//EXERCISE 11

	// var wg sync.WaitGroup
	// ch := make(chan int)
	// wg.Add(1)

	// go Exercise12(ch, &wg)

	// fmt.Println(<-ch)

	// wg.Wait()

	//EXERCISE 11

	var wg sync.WaitGroup
	ch := make(chan string)
	wg.Add(1)

	go Exercise13(ch, &wg)

	result := <-ch

	fmt.Println(result)

	wg.Wait()

}

func Exercise13(ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	num := 0
	fmt.Println("How many number in this string: ")
	fmt.Scan(&num)
	array := make([]string, num)
	for i:=0; i<num; i++ {
		fmt.Println("Enter ", i+1, " word: ")
		fmt.Scan(&array[i])
	}
	inputString := strings.Join(array, " ")
	words := strings.Fields(inputString)
	LongWord := ""
	for _, val := range words {
		if len(val) > len(LongWord) {
			LongWord = val
		}
	}

	ch <- LongWord

	close(ch)
}

func Exercise12(ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	num := 0
	fmt.Println("How many numbers do you want to enter? ")
	fmt.Scan(&num)
	slice := make([]int, num)
	sum := 0
	for i := 0; i < num; i++ {
		fmt.Println("Enter ", i+1, " number: ")
		fmt.Scan(&slice[i])
		if slice[i]%2 == 0 {
			sum += slice[i]
		}
	}

	ch <- sum

	close(ch)

}

func Exercise11(ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	num := 0
	fmt.Println("Enter number: ")
	fmt.Scan(&num)
	fib := make([]int, num)
	fib[0], fib[1] = 0, 1
	sum := 0
	for i := 2; i < num; i++ {
		fib[i] = fib[i-1] + fib[i-2]
		sum = sum + fib[i]
	}
	ch <- sum
	close(ch)
}

func Exercise10(ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	sum := 0
	for i := 1; i <= 10; i++ {
		sum += (i * i)
	}
	ch <- sum
	close(ch)
}

func Exercise9(ch chan []int, wg *sync.WaitGroup) {
	defer wg.Done()
	var num int
	fmt.Println("How many numbers do you want to enter? ")
	fmt.Scan(&num)
	slice := make([]int, num)
	for i := 0; i < num; i++ {
		fmt.Println("Enter a ", i+1, " number")
		fmt.Scan(&slice[i])
	}
	sort.Ints(slice)
	ch <- slice
	close(ch)
}

func Exercise8(ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	var num int
	fmt.Println("How many numbers do you want to enter? ")
	fmt.Scan(&num)
	slice := make([]int, num)
	sum := 0
	for i := 0; i < num; i++ {
		fmt.Println("Enter a ", i+1, " number")
		fmt.Scan(&slice[i])
		sum += slice[i]
	}

	ch <- sum / num
	close(ch)
}

func Exercise7(ch chan []int, wg *sync.WaitGroup) {
	defer wg.Done()
	var num int
	fmt.Println("How many numbers do you want to enter? ")
	fmt.Scan(&num)
	slice := make([]int, num)
	NewSlice := []int{}
	for i := 0; i < num; i++ {
		fmt.Println("Enter a ", i+1, " number")
		fmt.Scan(&slice[i])
		NewSlice = append(NewSlice, slice[i]*2)
	}

	ch <- NewSlice
	close(ch)
}

func Exercise6(ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("Enter word: ")
	var word string
	fmt.Scan(&word)

	l := len(word)
	NewWord := make([]rune, l)

	for i := l - 1; i >= 0; i-- {
		NewWord[l-1-i] = rune(word[i])
	}

	ch <- string(NewWord)

	close(ch)
}

func Exercise5(ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	var num int
	fmt.Println("How many numbers do you want to enter? ")
	fmt.Scan(&num)
	slice := make([]int, num)
	for i := 0; i < num; i++ {
		fmt.Println("Enter a ", i+1, " number")
		fmt.Scan(&slice[i])
	}
	sort.Ints(slice)
	max := slice[num-1]
	ch <- max
	close(ch)
}

func Exercise4(ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	var num int
	fmt.Println("Enter number: ")
	fmt.Scan(&num)
	sum := 0
	for i := 2; i < num; i++ {
		if num%2 == 0 {
			sum++
		}
	}
	if sum == 0 {
		ch <- "Prime Number"
	} else {
		ch <- "Not prime number"
	}

	close(ch)
}

func Exercise3(ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	var num int
	fact := 1
	fmt.Println("Enter number: ")
	fmt.Scan(&num)
	for i := 1; i <= num; i++ {
		fact *= i
	}
	ch <- fact
	close(ch)
}

func Exercise2(ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	counter := 0
	for i := 0; i <= 100; i++ {
		counter += i
	}

	ch <- counter
	close(ch)
}

func Exercise1(ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 1; i <= 10; i++ {
		ch <- i
	}
	close(ch)
}
