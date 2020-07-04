package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type AllRepositries []struct {
	Name string `json:"name"`
}

type Repositry []struct {
	Path        string      `json:"path"`
	URL         string      `json:"url"`
	DownloadURL interface{} `json:"download_url"`
	Type        string      `json:"type"`
}

var dirPerm os.FileMode = 0755
var dir = "dir"
var file = "file"

var errStatusNotOK = errors.New("Status Not 200")
var unexpectedErr = errors.New("Unexpected error occurred")
var nl = "\n"
var dl = "\n\n"

func main() {
	if err := run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

func run(args []string) error {
	switch {
	case len(args) == 1:
		showUsage()
	case len(args) > 2:
		showUsage()
	default:
		fmt.Println("Start importing...")
		userName := args[1]
		allRepoNames, err := getAllRepos(userName)
		if err != nil {
			return err
		}
		for _, repoName := range allRepoNames {
			url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/", userName, repoName)
			if err := saveRepoInLocal(url, repoName); err != nil {
				return err
			}
		}
		fmt.Println("Finished importing!!")
	}
	return nil
}

func saveRepoInLocal(url, repoName string) error {
	if url == "" || repoName == "" {
		return unexpectedErr
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errStatusNotOK
	}

	var repo = Repositry{}
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		return err
	}
	for _, repoContent := range repo {
		switch repoContent.Type {
		case dir:
			path := filepath.Join(repoName, repoContent.Path)
			if err := os.Mkdir(path, dirPerm); err != nil {
				fmt.Println(err)
				continue
			}
			if err := saveRepoInLocal(repoContent.URL, repoName); err != nil {
				return err
			}
		case file:
			if err := processFile(repoName, repoContent.Path, repoContent.DownloadURL.(string)); err != nil {
				return err
			}
		}
	}
	return nil
}

func processFile(repoName, filePath, downloadURL string) error {
	if repoName == "" || filePath == "" || downloadURL == "" {
		return unexpectedErr
	}
	req, err := http.NewRequest(http.MethodGet, downloadURL, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errStatusNotOK
	}

	path := filepath.Join(repoName, filePath)
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, resp.Body); err != nil {
		return err
	}
	return nil
}

func getAllRepos(userName string) ([]string, error) {
	if userName == "" {
		return nil, unexpectedErr
	}
	url := fmt.Sprintf("https://api.github.com/users/%s/repos", userName)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var allRepos = AllRepositries{}

	if err := json.NewDecoder(resp.Body).Decode(&allRepos); err != nil {
		return nil, err
	}
	var allRepoNames = []string{}
	for _, repo := range allRepos {
		allRepoNames = append(allRepoNames, repo.Name)
		if err := os.Mkdir(repo.Name, dirPerm); err != nil {
			fmt.Println(err)
			continue
		}
	}
	return allRepoNames, nil
}

func showUsage() {
	fmt.Println("This alows you to git clone all the public repositries of a user\n")
	fmt.Println("Usage:\n   gituser [github username]")
	fmt.Println("Description:\n   By executing git user [github username],\n   you can git clone all the public repositries of the user to your working directory.\n")
}
