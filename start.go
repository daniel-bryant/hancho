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

  var reply Configuration
  err = client.Call("ProxyManager.AddServices", config, &reply)
  if err != nil {
    log.Fatal("ProxyManager.AddServices error:", err)
  }

  *config = reply
  for serviceName, settings := range config.Services {
    fmt.Printf("%s.%s.localhost proxied to port %s\n", serviceName, config.Name, settings.Port)
  }
  fmt.Println()
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

  var reply ConfigMap
  err = client.Call("ProxyManager.GetServices", &Empty{}, &reply)
  if err != nil {
    log.Fatal("client.Call: ", err)
  }

  if len(reply) == 0 {
    fmt.Println("No projects registered")
    return
  }

  for projectName, config := range reply {
    for serviceName, settings := range config.Services {
      fmt.Printf("%s.%s.localhost proxied to port %s\n", serviceName, projectName, settings.Port)
    }
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

  var reply Configuration
  err = client.Call("ProxyManager.RemoveServices", config.Name, &reply)
  if err != nil {
    log.Fatal("ProxyManager.RemoveServices error:", err)
  }

  fmt.Println()
  for serviceName, settings := range reply.Services {
    fmt.Printf("Removed %s.%s.localhost proxy to port %s\n", serviceName, reply.Name, settings.Port)
  }
}
