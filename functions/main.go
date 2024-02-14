package main

import "fmt"

func main() {
fmt.Println(otherFunc(2,3,4))
}

func otherFunc(a,b,c int) (x, y int, z string) {
	x = a-1
	y = b+1
	
	return
}