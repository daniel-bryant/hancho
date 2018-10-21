package main

import (
  "path/filepath"
)

func startServices(config *Config) {
  for name, service := range config.Services {
    execCommand("docker-compose", "--file", composeFilepath(name, service), "up", "--detach")
  }
}

func stopServices(config *Config) {
  for name, service := range config.Services {
    execCommand("docker-compose", "--file", composeFilepath(name, service), "down")
  }
}

func composeFilepath(name string, service ServiceConfig) (string) {
  const composefile = "docker-compose.yml"

  if len(service.Localpath) != 0 {
    return filepath.Join(service.Localpath, composefile)
  }

  return filepath.Join(".hancho", name, composefile)
}
