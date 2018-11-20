package hancho

import (
  "fmt"
  "os"
  "path/filepath"

  "gopkg.in/src-d/go-git.v4"
)

func pullRepositories(config *Configuration) {
  gitDir := createDir(".hancho", "git")
  for name, settings := range config.Services {
    fmt.Printf("\nService: %s - %s\n", name, settings.GitUrl)
    pullRepository(filepath.Join(gitDir, name), settings.GitUrl)
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

func createDir(parent, name string) string {
  dir := filepath.Join(parent, name)
  err := os.MkdirAll(dir, os.ModePerm)
  checkError(err)

  return dir
}
