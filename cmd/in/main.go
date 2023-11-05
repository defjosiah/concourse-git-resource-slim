package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type InRequest struct {
	Source  map[string]interface{} `json:"source"`
	Version map[string]string      `json:"version"`
	Params  map[string]string      `json:"params"`
}

func main() {
	// Read from stdin
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read from stdin: %s\n", err)
		os.Exit(1)
	}

	// Unmarshal JSON input
	var request InRequest
	if err := json.Unmarshal(input, &request); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to unmarshal JSON: %s\n", err)
		os.Exit(1)
	}

	// Placeholder for in logic
	fmt.Println("In logic not implemented")
}
