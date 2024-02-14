package closurefunction

func Calculator() func(int, int) int {
	return func(x, y int) int {
		return x + y
	}

}
