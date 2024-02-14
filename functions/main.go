package main

import (
	"fmt"
	"test/factorial"
)

func main() {
	// fmt.Println(otherFunc(2,3,4))

	// add := closurefunction.Calculator()

	// result := add(2, 3)

	// fmt.Println(result)

	a := factorial.Factorial(5)
	fmt.Println(a)

}

func otherFunc(a, b, c int) (x, y int, z string) {
	x = a - 1
	y = b + 1

	return
}
