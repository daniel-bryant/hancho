package main

import (
  "path/filepath"
)

func startServices(config *Config) {
  for name, service := range config.Services {
    fp := composeFilepath(name, service)
    cmd := progressCommand("docker-compose", "--file", fp, "up", "--detach")
    cmd.EnvAppend("PORT", service.Port)
    cmd.Wait()
  }
}

func stopServices(config *Config) {
  for name, service := range config.Services {
    fp := composeFilepath(name, service)
    cmd := progressCommand("docker-compose", "--file", fp, "down")
    cmd.EnvAppend("PORT", service.Port)
    cmd.Wait()
  }
}

func composeFilepath(name string, service ServiceConfig) (string) {
  const composefile = "docker-compose.yml"

  if len(service.Localpath) != 0 {
    return filepath.Join(service.Localpath, composefile)
  }

  return filepath.Join(".hancho", name, composefile)
}
