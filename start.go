package main

func handleStartCommand() {
  config := Configuration()

  pullRepositories(config)
  startServices(config)
}

func handleStopCommand() {
  stopServices(Configuration())
}
