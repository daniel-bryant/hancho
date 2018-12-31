package hancho

import (
  "fmt"
  "os"
  "path/filepath"
)

func pullRepositories(config *Configuration) {
  gitDir := createDir(".hancho", "git")

  for name, settings := range config.Services {
    fmt.Printf("hancho: Updating service '%s'\n", name)

    dir := filepath.Join(gitDir, name)
    branch := settings.Branch
    repo := settings.GitUrl

    if len(branch) == 0 {
      branch = "master"
    }

    if _, err := os.Stat(dir); os.IsNotExist(err) {
      cmd := progressCommand("git", "clone", "--progress", "--single-branch", "--branch", branch, repo, dir)
      cmd.Wait()
    } else {
      cmd := progressCommand("git", "fetch", "--progress", "origin", branch)
      cmd.SetDir(dir)
      cmd.Wait()

      cmd = progressCommand("git", "reset", "--hard", "origin/" + branch)
      cmd.SetDir(dir)
      cmd.Wait()
    }
  }

  fmt.Println()
}

func createDir(parent, name string) string {
  dir := filepath.Join(parent, name)
  err := os.MkdirAll(dir, os.ModePerm)
  checkError(err)

  return dir
}
