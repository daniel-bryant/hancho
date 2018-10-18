package main

func handleStartCommand() {
  config := getConfig()
  pullRepositories(config)
  startServices(config)
}

func handleStopCommand() {
  config := getConfig()
  stopServices(config)
}
