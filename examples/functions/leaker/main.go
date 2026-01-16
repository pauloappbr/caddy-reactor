package main

import (
	"fmt"
	"os"
)

func main() {
	os.Stdout.WriteString("💧 Trying to allocate 100MB of RAM...\n")

	size := 100 * 1024 * 1024

	data := make([]byte, size)

	data[0] = 1
	data[size-1] = 1

	fmt.Printf("✅ Success! I allocated %d bytes.\n", len(data))
}
