package git

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

func Clone(url, tag string) (string, error) {
	// Clone the given repository to the given directory
	log.Printf("git clone %s", url)
	splitted := strings.Split(url, "/")
	repoName := splitted[len(splitted)-1]
	directory := fmt.Sprintf("/tmp/%s", repoName)
	log.Printf("name %s", repoName)
	log.Printf("tag %s", tag)

	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL: url,
	})
	if err != nil && err != git.ErrRepositoryAlreadyExists {
		return "", fmt.Errorf("plain clone: %s", err)
	}

	if err == git.ErrRepositoryAlreadyExists {
		r, err = git.PlainOpen(directory)
		if err != nil {
			return "", fmt.Errorf("plain open: %s", err)
		}
	}

	log.Print("git show-ref --head HEAD")
	ref, err := r.Head()
	if err != nil {
		return "", fmt.Errorf("head: %s", err)
	}

	fmt.Println(ref.Hash())

	w, err := r.Worktree()
	if err != nil {
		return "", fmt.Errorf("worktree: %s", err)
	}

	err = r.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return "", fmt.Errorf("fetch: %s", err)
	}

	// ... checking out to commit
	log.Printf("git checkout %s", tag)
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName("websockets"),
	})

	log.Print("git show-ref --head HEAD")
	ref, err = r.Head()
	if err != nil {
		return "", fmt.Errorf("head: %s", err)
	}
	fmt.Println(ref.Hash())

	return directory, nil
}
