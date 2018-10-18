package main

import (
  "fmt"
  "os/exec"
  "path/filepath"
)

func startServices(config *Config) {
  for _, service := range config.Services {
    composeFilepath := filepath.Join(service.Localpath, "docker-compose.yml")
    cmd := exec.Command("docker-compose", "--file", composeFilepath, "up", "-d")
    out, err := cmd.CombinedOutput()
    checkError(err)
    fmt.Printf("docker-compose up output:\n%s\n", string(out))
  }
}

func stopServices(config *Config) {
  for _, service := range config.Services {
    composeFilepath := filepath.Join(service.Localpath, "docker-compose.yml")
    cmd := exec.Command("docker-compose", "--file", composeFilepath, "down")
    out, err := cmd.CombinedOutput()
    checkError(err)
    fmt.Printf("docker-compose down output:\n%s\n", string(out))
  }
}
