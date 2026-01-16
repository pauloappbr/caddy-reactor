package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Order struct {
	ID    string  `json:"id"`
	Value float64 `json:"value"`
}

func main() {
	var order Order
	decoder := json.NewDecoder(os.Stdin)
	if err := decoder.Decode(&order); err != nil {
		os.Stdout.WriteString(`{"error": "json invalid"}`)
		return
	}

	tax := order.Value * 0.15
	total := order.Value + tax

	fmt.Printf(`{"order_id": "%s", "total_final": %.2f}`, order.ID, total)
}
