package hancho

import (
  "log"
)

const (
  ProxyManagerPort = ":1234"
)

type Empty struct {}

type ProxyManager struct {
  Services ServiceMap
}

func NewProxyManager() (*ProxyManager) {
  return &ProxyManager{make(ServiceMap)}
}

func (p *ProxyManager) AddServices(services ServiceMap, reply *int) error {
  log.Printf("Adding %d services\n", len(services))
  for name, settings := range services {
    p.Services[name] = settings
  }
  *reply = len(services)
  return nil
}

func (p *ProxyManager) GetServices(empty *Empty, reply *ServiceMap) error {
  log.Printf("Returning %d services\n", len(p.Services))
  *reply = p.Services
  return nil
}
