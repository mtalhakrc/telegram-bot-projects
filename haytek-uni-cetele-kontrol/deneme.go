package main

import (
	"fmt"
	"math"
	"strconv"
)

func main() {
	a := math.NaN()
	str := strconv.Itoa(int(a))
	fmt.Println(int(a))
	fmt.Println(str)
}
