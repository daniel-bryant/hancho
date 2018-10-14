package main

import (
  "log"
  "io/ioutil"

  "gopkg.in/yaml.v2"
)

type Config struct {
  Services map[string]ServiceConfig
}

type ServiceConfig struct {
  Language string
  Protocol string
  Giturl string
  Localpath string
}

func getConfig() (*Config) {
  type ProjectConfig struct {
    Services map[string]struct {
      Language string
      Protocol string
      Giturl string
    }
  }

  type LocalConfig struct {
    Services map[string]struct {
      Localpath string
    }
  }

  projectConfig := ProjectConfig{}
  localConfig := LocalConfig{}

  readConfig("config.yml", &projectConfig)
  readConfig(".config.local.yml", &localConfig)

  services := make(map[string]ServiceConfig)
  for name, sc := range projectConfig.Services {
    localpath := localConfig.Services[name].Localpath
    services[name] = ServiceConfig{sc.Language, sc.Protocol, sc.Giturl, localpath}
  }

  config := Config{services}
  return &config
}

func readConfig(filename string, out interface{}) {
  data, err := ioutil.ReadFile(filename)
  if err != nil {
    log.Fatal(err)
  }

  err = yaml.Unmarshal(data, out)
  if err != nil {
    log.Fatalf("error: %v", err)
  }
}
