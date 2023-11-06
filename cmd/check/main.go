package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"time"
)

type CheckRequest struct {
	Source  Source   `json:"source"`
	Version *Version `json:"version,omitempty"`
}

type Source struct {
	Branch    string     `json:"branch"`
	Paths     []string   `json:"paths"`
	Repo      Repository `json:"repo"`
	AuthToken string     `json:"auth-token"`
}

type Repository struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

type Version struct {
	Ref string `json:"ref"`
}

type Commit struct {
	Sha    string `json:"sha"`
	Commit struct {
		Committer struct {
			Name string `json:"name"`
			Date string `json:"date"`
		} `json:"committer"`
		Message string `json:"message"`
	} `json:"commit"`
}

func fetchCommits(source Source, path string) ([]Commit, error) {
	// Create a new HTTP client
	client := &http.Client{}

	// Construct the API URL
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits?sha=%s&path=%s", source.Repo.Owner, source.Repo.Name, source.Branch, path)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	// Add headers to the request
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", "Bearer "+source.AuthToken)
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	// Perform the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API responded with status: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Decode the JSON response into the Commit slice
	var commits []Commit
	if err := json.Unmarshal(body, &commits); err != nil {
		return nil, err
	}

	return commits, nil
}

func main() {
	// Read from stdin
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read from stdin: %s\n", err)
		os.Exit(1)
	}

	// Unmarshal JSON input
	var request CheckRequest
	if err := json.Unmarshal(input, &request); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to unmarshal JSON: %s\n", err)
		os.Exit(1)
	}

	allCommits := []Commit{}
	// Fetch the commits using the new function
	for _, path := range request.Source.Paths {
		commits, err := fetchCommits(request.Source, path)
		allCommits = append(allCommits, commits...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching commits: %s\n", err)
			os.Exit(1)
		}
	}
	// Sort all commits by date
	sort.Slice(allCommits, func(i, j int) bool {
		timeI, errI := time.Parse(time.RFC3339, allCommits[i].Commit.Committer.Date)
		timeJ, errJ := time.Parse(time.RFC3339, allCommits[j].Commit.Committer.Date)
		if errI != nil || errJ != nil {
			fmt.Fprintf(os.Stderr, "Error parsing commit dates: %s %s\n", errI, errJ)
			os.Exit(1)
		}
		return timeI.After(timeJ)
	})

	if len(allCommits) == 0 {
		fmt.Fprintf(os.Stderr, "No new versions\n")
		os.Exit(0)
		return
	}

	// version when empty
	if request.Version == nil {
		versions := []Version{{Ref: allCommits[0].Sha}}
		fmt.Fprintf(os.Stderr, "No incoming version, picking latest\nAuthor: %s\nDate: %s\nMessage: %s\n",
			allCommits[0].Commit.Committer.Name,
			allCommits[0].Commit.Committer.Date,
			allCommits[0].Commit.Message)
		output, err := json.Marshal(versions)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to marshal JSON: %s\n", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
		return
	}

	afterRef := request.Version.Ref
	afterRefVersions := []Version{}
	fmt.Fprintf(os.Stderr, "Incoming version, selecting commits after: %s\n", afterRef)
	for _, commit := range allCommits {
		if commit.Sha == afterRef {
			break
		}
		fmt.Fprintf(os.Stderr, "\nAuthor: %s\nDate: %s\nMessage: %s\n",
			commit.Commit.Committer.Name,
			commit.Commit.Committer.Date,
			commit.Commit.Message)
		afterRefVersions = append(afterRefVersions, Version{Ref: commit.Sha})
	}
	if len(afterRefVersions) == 0 {
		fmt.Fprintf(os.Stderr, "No new versions\n")
		os.Exit(0)
		return
	}
	output, err := json.Marshal(afterRefVersions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal JSON: %s\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "%s", output)
}
