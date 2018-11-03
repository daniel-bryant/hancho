package main

import (
  "fmt"
  "os"
  "path/filepath"

  "gopkg.in/src-d/go-git.v4"
)

func pullRepositories(config *Config) {
  gitDir := createDir(".hancho", "git")
  for _, service := range config.Services {
    fmt.Printf("\nService: %s - %s\n", service.Name, service.GitUrl)
    pullRepository(filepath.Join(gitDir, service.Name), service.GitUrl)
  }
}

func pullRepository(directory, url string) {
    repository, err := git.PlainOpen(directory)

    if err == git.ErrRepositoryNotExists {
      _, err := git.PlainClone(directory, false, &git.CloneOptions{
        URL:               url,
        Progress:          os.Stdout,
        RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
      })

      checkError(err)
      return
    }

    checkError(err)

    worktree, err := repository.Worktree()

    err = worktree.Pull(&git.PullOptions{
      RemoteName: "origin",
      Progress:   os.Stdout,
    })

    if err == git.NoErrAlreadyUpToDate {
      fmt.Println("Already up to date.")
      return
    }

    checkError(err)
}
