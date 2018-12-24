package hancho

import (
  "log"
  "strconv"
)

const (
  ProxyManagerPort = ":1234"
  startingPort = 5000
)

type Empty struct {}

type ProxyManager struct {
  Services ServiceMap
  NextAvailablePort int
}

func NewProxyManager() (*ProxyManager) {
  return &ProxyManager{make(ServiceMap), startingPort}
}

func (p *ProxyManager) AddServices(services ServiceMap, reply *ServiceMap) error {
  log.Printf("Adding %d services\n", len(services))
  for name, settings := range services {
    if len(settings.Port) == 0 {
      settings.Port = strconv.Itoa(p.NextAvailablePort)
      p.NextAvailablePort += 1
    }
    services[name] = settings
    p.Services[name] = settings
  }
  *reply = services
  return nil
}

func (p *ProxyManager) RemoveServices(services ServiceMap, reply *ServiceMap) error {
  log.Printf("Removing %d services\n", len(services))
  for name, _ := range services {
    delete(p.Services, name)
  }
  *reply = services
  return nil
}

func (p *ProxyManager) GetServices(empty *Empty, reply *ServiceMap) error {
  log.Printf("Returning %d services\n", len(p.Services))
  *reply = p.Services
  return nil
}
