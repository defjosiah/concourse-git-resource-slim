package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
)

type OutRequest struct {
    Source map[string]interface{} `json:"source"`
    Params map[string]string      `json:"params"`
}

func main() {
    // Read from stdin
    input, err := ioutil.ReadAll(os.Stdin)
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

    // Placeholder for out logic
    fmt.Println("Out logic not implemented")
}
