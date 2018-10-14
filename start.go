package main

func handleStartCommand() {
  config := getConfig()
  pullRepositories(config)
}
