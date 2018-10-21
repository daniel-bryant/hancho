package main

import (
  "path/filepath"
)

func startServices(config *Config) {
  for _, service := range config.Services {
    composeFilepath := filepath.Join(service.Localpath, "docker-compose.yml")
    execCommand("docker-compose", "--file", composeFilepath, "up", "--detach")
  }
}

func stopServices(config *Config) {
  for _, service := range config.Services {
    composeFilepath := filepath.Join(service.Localpath, "docker-compose.yml")
    execCommand("docker-compose", "--file", composeFilepath, "down")
  }
}
