package hancho

import (
  "log"
)

const (
  ProxyManagerPort = ":1234"
)

type ProxyManager struct {
  Services []Service
}

func (p *ProxyManager) AddServices(services []Service, reply *int) error {
  log.Printf("Adding %d services\n", len(services))
  p.Services = append(p.Services, services...)
  *reply = len(services)
  return nil
}
