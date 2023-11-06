package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type OutRequest struct {
	Params  map[string]string `json:"params"`
	Source  Source            `json:"source"`
	Version *Version          `json:"version,omitempty"`
}

type Repository struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

type Source struct {
	Branch    string     `json:"branch"`
	Paths     []string   `json:"paths"`
	Repo      Repository `json:"repo"`
	AuthToken string     `json:"auth-token"`
}

type Version struct {
	Ref string `json:"ref"`
}

func main() {
	// Read from stdin
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read from stdin: %s\n", err)
		os.Exit(1)
	}

	// Unmarshal JSON input
	var request OutRequest
	if err := json.Unmarshal(input, &request); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to unmarshal JSON: %s\n", err)
		os.Exit(1)
	}

	// Write the version to stdout
	output, err := json.Marshal(request.Version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal JSON: %s\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "%s", output)
}
