package main

import (
  "os"
)

func handleStartCommand() {
  os.Mkdir(".hancho", os.ModePerm)

  config := Configuration()

  pullRepositories(config)
  startServices(config)
  startProxyServer(config)
  stopServices(config)
}

func handleStopCommand() {
  stopServices(Configuration())
}
