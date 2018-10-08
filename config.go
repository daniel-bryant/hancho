package main

import (
  "log"
  "io/ioutil"

  "gopkg.in/yaml.v2"
)

type Config struct {
  Version int
  Services []struct {
    Name string
    Language string
    Protocol string
    Giturl string
  }
}

func readConfig() (*Config) {
  config := Config{}

  data, err := ioutil.ReadFile("config.yml")
  if err != nil {
    log.Fatal(err)
  }

  err = yaml.Unmarshal(data, &config)
  if err != nil {
    log.Fatalf("error: %v", err)
  }

  return &config
}
