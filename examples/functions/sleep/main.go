package main

import (
	"fmt"
	"os"
)

func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

func main() {
	os.Stdout.WriteString("🔥 Starting Fibonacci of Death (fib 45)...\n")

	result := fib(45)

	fmt.Printf("😱 I survived! Result: %d\n", result)
}
