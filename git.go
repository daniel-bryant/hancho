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

    exists := true
    if _, err := os.Stat(dir); os.IsNotExist(err) {
      exists = false
    }

    if exists {
      cmd := progressCommand("git", "pull")
      cmd.SetDir(dir)
      cmd.Wait()
    } else {
      cmd := progressCommand("git", "clone", "--progress", settings.GitUrl, dir)
      cmd.Wait()
    }

    if len(settings.Branch) != 0 {
      cmd := progressCommand("git", "checkout", settings.Branch)
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
