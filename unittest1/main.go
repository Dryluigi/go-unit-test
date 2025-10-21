package main

import (
	"fmt"
	"unit-test-demo/unittest1/lib"
)

func main() {
	res, _ := lib.FormatDateLong("2025-10-18")

	fmt.Println(res)
}