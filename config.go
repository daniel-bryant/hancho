package hancho

import (
  "io/ioutil"
  "log"

  "gopkg.in/yaml.v2"
)

type Config struct {
  Services []Service
}

type Service struct {
  Name string
  ServiceSettings
}

type ConfigFile struct {
  Services map[string] ServiceSettings
}

type ServiceSettings struct {
  Port string
  GitUrl string
  LocalPath string
}

func Configuration() (*Config) {
  configFile := ConfigFile{}
  localConfigFile := ConfigFile{}

  unmarshalConfig("config.yml", &configFile)
  unmarshalConfig(".config.yml", &localConfigFile)

  for name, localSettings := range localConfigFile.Services {
    if remoteSettings, exists := configFile.Services[name]; exists {
      port := remoteSettings.Port
      if len(localSettings.Port) != 0 {
        port = localSettings.Port
      }

      gitUrl := remoteSettings.GitUrl
      if len(localSettings.GitUrl) != 0 {
        gitUrl = localSettings.GitUrl
      }

      localPath := remoteSettings.LocalPath
      if len(localSettings.LocalPath) != 0 {
        localPath = localSettings.LocalPath
      }

      configFile.Services[name] = ServiceSettings{port, gitUrl, localPath}
    } else {
      configFile.Services[name] = localSettings
    }
  }

  var services []Service
  for name, serviceSettings := range configFile.Services {
    services = append(services, Service{name, serviceSettings})
  }

  return &Config{services}
}

func unmarshalConfig(filename string, out *ConfigFile) {
  data, err := ioutil.ReadFile(filename)
  if err != nil {
    log.Fatal(err)
  }

  err = yaml.Unmarshal(data, out)
  if err != nil {
    log.Fatalf("error: %v", err)
  }
}
