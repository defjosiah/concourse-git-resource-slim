package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type InRequest struct {
	Source  *Source           `json:"source"`
	Version *Version          `json:"version"`
	Params  map[string]string `json:"params"`
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

// DownloadFile will download a url to a local file.
func DownloadFile(source *Source, ref string, archiveFile *os.File) error {
	client := &http.Client{
		Timeout: time.Minute * 10,
	}

	apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/tarball/%s", source.Repo.Owner, source.Repo.Name, ref)
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/vnd.github.v3.raw")
	req.Header.Set("Authorization", "Bearer "+source.AuthToken)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	defer archiveFile.Close()
	_, err = io.Copy(archiveFile, resp.Body)
	return err
}

// Untar will untar a specified tarball to a destination directory.
// Just use tar cli instead of attempting to do it in golang.
// tar is really good at that, go is less good at that.
func Untar(tarball, destination string) error {
	// strip the first component so we don't have an extra nested folder
	cmd := exec.Command("tar", "-xzf", tarball, "--strip-components=1", "-C", destination)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func writeGitInformationFiles(source *Source, ref string, destinationPath string) {
	client := &http.Client{
		Timeout: time.Minute * 10,
	}
	apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s", source.Repo.Owner, source.Repo.Name, ref)
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create request: %s\n", err)
		os.Exit(1)
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", "Bearer "+source.AuthToken)
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to perform request: %s\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Received non-200 status code: %d\n", resp.StatusCode)
		os.Exit(1)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read body: %s\n", err)
		os.Exit(1)
	}

	var commit Commit
	if err := json.Unmarshal(body, &commit); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to unmarshal JSON: %s\n", err)
		os.Exit(1)
	}

	// Function to write content to a file within the .git directory
	writeFile := func(filename, content string) {
		filePath := filepath.Join(destinationPath, ".git", filename)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write %s file: %s\n", filename, err)
			os.Exit(1)
		}
	}

	// Create the .git directory if it does not exist
	gitDir := filepath.Join(destinationPath, ".git")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create .git directory: %s\n", err)
		os.Exit(1)
	}

	writeFile("committer", commit.Commit.Committer.Name)
	writeFile("ref", ref)
	writeFile("short_ref", ref[:7])
	writeFile("commit_message", commit.Commit.Message)
	writeFile("commit_timestamp", commit.Commit.Committer.Date)
}

func main() {
	// get first argument from command line
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <destination>\n", os.Args[0])
		os.Exit(1)
	}
	// The in script is passed a destination directory as command line argument $1,
	// and is given on stdin the configured source and a precise version of the resource to fetch.
	destinationPath := os.Args[1]

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read from stdin: %s\n", err)
		os.Exit(1)
	}
	var request InRequest
	if err := json.Unmarshal(input, &request); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to unmarshal JSON: %s\n", err)
		os.Exit(1)
	}

	archiveFile, err := os.CreateTemp("", "git-*.tar.gz")
	fmt.Fprintf(os.Stderr, "Writing to archive: %s\n", archiveFile.Name())
	requestedVersion := request.Version.Ref

	// write ".git" files to the destination
	writeGitInformationFiles(request.Source, requestedVersion, destinationPath)

	if err := DownloadFile(request.Source, requestedVersion, archiveFile); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to download file: %s\n", err)
		os.Exit(1)
	}

	if err := Untar(archiveFile.Name(), destinationPath); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to untar file: %s\n", err)
		os.Exit(1)
	}
	output, err := json.Marshal(request.Version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal JSON: %s\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "%s", output)
}
