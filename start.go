package hancho

import (
  "log"
  "net/rpc"
)

func handleStartCommand() {
  config := ReadConfiguration()

  pullRepositories(config)
  registerProxies(config)
  startServices(config)
}

func registerProxies(config *Configuration) {
  serverAddress := "localhost"
  client, err := rpc.DialHTTP("tcp", serverAddress + ProxyManagerPort)
  if err != nil {
    log.Fatal("dialing:", err)
  }

  var reply int
  err = client.Call("ProxyManager.AddServices", config.Services, &reply)
  if err != nil {
    log.Fatal("ProxyManager.AddServices error:", err)
  }

  log.Printf("Added %d servers\n", reply)
}

func handleStopCommand() {
  stopServices(ReadConfiguration())
}
