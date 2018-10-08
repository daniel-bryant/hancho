package main

func handleStartCommand() {
  config := readConfig()
  pullRepositories(config)
}
