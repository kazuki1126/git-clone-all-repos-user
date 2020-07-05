package main

import (
	"fmt"
	git "go-excercises/gitcloneCLI"
	"log"
	"os"
)

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
		fmt.Println("Invalid number of arguments")
		showUsage()
	default:
		fmt.Println("Started importing...")
		userName := args[1]
		allRepoNames, err := git.GetAllRepos(userName)
		if err != nil {
			return err
		}
		for _, repoName := range allRepoNames {
			url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/", userName, repoName)
			if err := git.CreateRepoInLocal(url, repoName); err != nil {
				return err
			}
		}
		fmt.Println("Finished importing!!")
	}
	return nil
}

func showUsage() {
	fmt.Println("This alows you to git clone all the public repositries of a user")
	fmt.Println("\nUsage:\n   gituser [github username]")
	fmt.Println("Description:\n   By executing git user [github username],\n   you can git clone all the public repositries of the user to your working directory.")
}
