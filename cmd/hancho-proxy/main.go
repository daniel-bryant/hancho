/*
 * ReverseProxy builder and Transport implementation
 *
 * https://medium.com/@jnewmano/grpc-postman-173b62a64341
 * https://github.com/jnewmano/grpc-json-proxy/blob/master/proxy.go
 *
 * https://golang.org/pkg/net/http/httputil/#example_ReverseProxy
 */

package main

import (
  "crypto/tls"
  "fmt"
  "log"
  "net"
  "net/http"
  "net/http/httptest"
  "net/http/httputil"
  "net/rpc"
  "net/url"
  "strings"
  "time"

  "github.com/daniel-bryant/hancho"
  "golang.org/x/net/http2"
  "golang.org/x/net/http2/h2c"
)

const (
  addr = ":8080"
  defaultClientTimeout = time.Second * 60
)

type Server struct {
  Subdomain string
  Port string
}

type Transport struct {
  HTTPClient    *http.Client
  H2Client      *http.Client
  H2NoTLSClient *http.Client
}

func (t Transport) RoundTrip(r *http.Request) (*http.Response, error) {
  client := t.HTTPClient
  if isGRPC(r) {
    if isSecure(r) {
      client = t.H2Client
    } else {
      client = t.H2NoTLSClient
    }
  }

  // clear requestURI, set in call to director
  r.RequestURI = ""

  return client.Do(r)
}

func main() {
  // Create and register a ProxyManager
  pm := hancho.NewProxyManager()
  rpc.Register(pm)
  rpc.HandleHTTP()

  l, e := net.Listen("tcp", hancho.ProxyManagerPort)
  if e != nil {
    log.Fatal("listen error:", e)
  }
  go http.Serve(l, nil)

  // Not found server
  notFound := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Service not found.")
  }))
  defer notFound.Close()

  director := func(r *http.Request) {
    parts := strings.Split(r.Host, ".")
    u := notFound.URL

    // example parts array: [dashboard checkr localhost:80]
    if len(parts) == 3 {
      serviceName := parts[0]
      projectName := parts[1]
      config := pm.Configs[projectName]
      service := config.Services[serviceName]

      if len(service.Port) != 0 {
        u = "http://localhost:" + service.Port
      }
    }

    origin, err := url.Parse(u)
    if err != nil {
      log.Fatal(err)
    }

    r.URL.Host = origin.Host
    r.URL.Scheme = origin.Scheme
    r.Header.Add("X-Forwarded-Host", r.Host)
    r.Header.Add("X-Origin-Host", origin.Host)

    reqType := "HTTP"
    if isGRPC(r) { reqType = "GRPC" }
    log.Printf("Forwarding %s request from `%s` to `%s`\n", reqType, r.Host, origin.Host)
  }

  h2NoTLSClient := &http.Client{
    // Skip TLS dial
    Transport: &http2.Transport{
      AllowHTTP: true,
      DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) {
        return net.Dial(netw, addr)
      },
    },
    Timeout: defaultClientTimeout,
  }

  h2Client := &http.Client{
    Transport: &http2.Transport{},
    Timeout:   defaultClientTimeout,
  }

  client := &http.Client{
    Timeout: defaultClientTimeout,
  }

  transport := &Transport{
    HTTPClient:    client,
    H2Client:      h2Client,
    H2NoTLSClient: h2NoTLSClient,
  }

  proxy := &httputil.ReverseProxy{Director: director, Transport: transport}
  h1Handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    proxy.ServeHTTP(w, r)
  })

  h2s := &http2.Server{}
  h1s := &http.Server{
    Addr: addr,
    Handler: h2c.NewHandler(h1Handler, h2s),
  }

  log.Println("starting HTTP server on port", addr)
  log.Fatal(h1s.ListenAndServe())
}

func isGRPC(r *http.Request) bool {
  if r.Header.Get("Content-Type") == "application/grpc" {
    return true
  }

  return false
}

func isSecure(r *http.Request) bool {
  return false
}
