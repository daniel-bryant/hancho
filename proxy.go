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
type ConfigMap map[string] Configuration

type ProxyManager struct {
  Configs ConfigMap
  NextAvailablePort int
}

func NewProxyManager() (*ProxyManager) {
  return &ProxyManager{make(ConfigMap), startingPort}
}

func (p *ProxyManager) AddServices(config Configuration, reply *Configuration) error {
  for name, settings := range config.Services {
    if len(settings.Port) == 0 {
      settings.Port = strconv.Itoa(p.NextAvailablePort)
      config.Services[name] = settings
      p.NextAvailablePort += 1
    }
  }
  p.Configs[config.Name] = config
  *reply = config

  for name, settings := range config.Services {
    log.Printf("%s.%s.localhost proxied to port %s\n", name, config.Name, settings.Port)
  }

  return nil
}

func (p *ProxyManager) RemoveServices(name string, reply *Configuration) error {
  *reply = p.Configs[name]
  delete(p.Configs, name)

  log.Printf("Removed `%s` services\n", name)
  for serviceName, _ := range reply.Services {
    log.Printf("Removed %s.%s.localhost", serviceName, name)
  }

  return nil
}

func (p *ProxyManager) GetServices(empty *Empty, reply *ConfigMap) error {
  log.Printf("Returning %d config(s)\n", len(p.Configs))
  *reply = p.Configs
  return nil
}
