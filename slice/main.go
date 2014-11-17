package main

import "fmt"

// +gen slice:"Max"
type Price float64

func main() {
	prices := PriceSlice{12.34, 43.21}

	fmt.Println(prices)
}
