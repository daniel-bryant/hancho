package main

import (
  "os"
)

func handleStartCommand() {
  os.Mkdir(".hancho", os.ModePerm)

  config := getConfig()

  pullRepositories(config)
  startServices(config)
  startProxyServer(config)
  stopServices(config)
}

func handleStopCommand() {
  config := getConfig()
  stopServices(config)
}
