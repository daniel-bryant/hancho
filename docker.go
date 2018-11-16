package hancho

import (
  "path/filepath"
)

func startServices(config *Config) {
  for _, service := range config.Services {
    fp := composeFilepath(service)
    cmd := progressCommand("docker-compose", "--file", fp, "up", "--detach")
    cmd.EnvAppend("PORT", service.Port)
    cmd.Wait()
  }
}

func stopServices(config *Config) {
  for _, service := range config.Services {
    fp := composeFilepath(service)
    cmd := progressCommand("docker-compose", "--file", fp, "down")
    cmd.EnvAppend("PORT", service.Port)
    cmd.Wait()
  }
}

func composeFilepath(service Service) (string) {
  const composefile = "docker-compose.yml"

  if len(service.LocalPath) != 0 {
    return filepath.Join(service.LocalPath, composefile)
  }

  return filepath.Join(".hancho", service.Name, composefile)
}
