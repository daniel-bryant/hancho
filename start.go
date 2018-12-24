package hancho

import (
  "log"
  "fmt"
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

  var reply ServiceMap
  err = client.Call("ProxyManager.AddServices", config.Services, &reply)
  if err != nil {
    log.Fatal("ProxyManager.AddServices error:", err)
  }

  config.Services = reply
  log.Printf("Registered %d proxies\n\n", len(reply))
}

func handleProxiesCommand() {
  serverAddress := "localhost"
  client, err := rpc.DialHTTP("tcp", serverAddress + ProxyManagerPort)
  if err != nil {
    fmt.Println("Could not access the proxy server. Is it running?")
    fmt.Println("Start it with 'hancho-proxy'.")
    fmt.Println()

    log.Fatal("Error:", err)
  }

  var reply ServiceMap
  err = client.Call("ProxyManager.GetServices", &Empty{}, &reply)
  if err != nil {
    log.Fatal("client.Call: ", err)
  }

  if len(reply) == 0 {
    fmt.Println("No services registered")
    return
  }

  for name, settings := range reply {
    fmt.Printf("%s is proxied to port %s\n", name, settings.Port)
  }
}

func handleStopCommand() {
  // todo: should save this config somewhere instead of rereading
  // in case changes were made between starting and stopping services
  config := ReadConfiguration()

  stopServices(config)

  serverAddress := "localhost"
  client, err := rpc.DialHTTP("tcp", serverAddress + ProxyManagerPort)
  if err != nil {
    log.Fatal("dialing:", err)
  }

  var reply ServiceMap
  err = client.Call("ProxyManager.RemoveServices", config.Services, &reply)
  if err != nil {
    log.Fatal("ProxyManager.RemoveServices error:", err)
  }

  log.Printf("Unregistered %d proxies\n\n", len(reply))
}
