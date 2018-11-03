package main

func handleStartCommand() {
  config := Configuration()

  pullRepositories(config)
  startServices(config)
  startProxyServer(config)
  stopServices(config)
}

func handleStopCommand() {
  stopServices(Configuration())
}
