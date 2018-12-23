package hancho

import (
  "fmt"
  "path/filepath"
)

func startServices(config *Configuration) {
  for name, settings := range config.Services {
    fmt.Printf("hancho: Starting service '%s'\n", name)
    fp := composeFilepath(name, settings)
    cmd := progressCommand("docker-compose", "--file", fp, "up", "--detach")
    if len(settings.Port) != 0 {
      cmd.EnvAppend("PORT", settings.Port)
    }
    cmd.Wait()
  }
}

func stopServices(config *Configuration) {
  for name, settings := range config.Services {
    fmt.Printf("hancho: Stopping service '%s'\n", name)
    fp := composeFilepath(name, settings)
    cmd := progressCommand("docker-compose", "--file", fp, "down")
    if len(settings.Port) != 0 {
      cmd.EnvAppend("PORT", settings.Port)
    }
    cmd.Wait()
  }
}

func composeFilepath(name string, settings ServiceSettings) (string) {
  const composefile = "docker-compose.yml"

  if len(settings.LocalPath) != 0 {
    return filepath.Join(settings.LocalPath, composefile)
  }

  return filepath.Join(".hancho", "git", name, composefile)
}
