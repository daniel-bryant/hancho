package hancho

import (
  "io/ioutil"
  "log"

  "gopkg.in/yaml.v2"
)

type ServiceMap map[string] ServiceSettings

type ServiceSettings struct {
  Port string
  GitUrl string
  Branch string
  LocalPath string
}

type Configuration struct {
  Name string
  Services ServiceMap
}

func ReadConfiguration() (*Configuration) {
  config := Configuration{}
  localConfig := Configuration{}

  unmarshalConfig("config.yml", &config)
  unmarshalConfig(".config.yml", &localConfig)

  for name, localSettings := range localConfig.Services {
    if remoteSettings, exists := config.Services[name]; exists {
      port := remoteSettings.Port
      if len(localSettings.Port) != 0 {
        port = localSettings.Port
      }

      gitUrl := remoteSettings.GitUrl
      if len(localSettings.GitUrl) != 0 {
        gitUrl = localSettings.GitUrl
      }

      branch := remoteSettings.Branch
      if len(localSettings.Branch) != 0 {
        branch = localSettings.Branch
      }

      localPath := remoteSettings.LocalPath
      if len(localSettings.LocalPath) != 0 {
        localPath = localSettings.LocalPath
      }

      config.Services[name] = ServiceSettings{port, gitUrl, branch, localPath}
    } else {
      config.Services[name] = localSettings
    }
  }

  return &config
}

func unmarshalConfig(filename string, out *Configuration) {
  data, err := ioutil.ReadFile(filename)
  if err != nil {
    log.Fatal(err)
  }

  err = yaml.Unmarshal(data, out)
  if err != nil {
    log.Fatalf("error: %v", err)
  }
}
