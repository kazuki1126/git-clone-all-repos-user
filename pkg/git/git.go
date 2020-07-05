package git

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func CreateRepoInLocal(url, repoName string) error {
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
			if err := createRepoInLocal(repoContent.URL, repoName); err != nil {
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

func GetAllRepos(userName string) ([]string, error) {
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
