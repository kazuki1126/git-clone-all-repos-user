package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sync"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	token = flag.String("t", "", "your github token")
	org   = flag.String("org", "", "organization to git clone repositries from")
)

// Is there any alternative for authorization other than the token
func main() {
	flag.Parse()
	if *token == "" || *org == "" {
		fmt.Println("Please specify your github token and organiation to clone repositries from.")
		os.Exit(1)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	repos, _, err := client.Repositories.ListByOrg(ctx, *org, nil)
	if err != nil {
		log.Fatal(err)
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	var wg sync.WaitGroup
	for _, repo := range repos {
		wg.Add(1)
		go func(repo *github.Repository, wg *sync.WaitGroup) {
			cmd := exec.Command("git", "clone", *repo.CloneURL)
			if err := cmd.Run(); err != nil {
				fmt.Printf("Error cloning %s Error:%s\n", *repo.CloneURL, err)
				return
			} else {
				fmt.Printf("git cloned %s\n", *repo.CloneURL)
				wg.Done()
			}
		}(repo, &wg)
	}
	wg.Wait()
}
