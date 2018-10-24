package main

import (
  "bufio"
  "fmt"
  "os"
  "os/signal"
  "path/filepath"
  "syscall"
  "text/template"
)

type NginxConf struct {
  AccessLog string
}

type ServerConf struct {
  Host string
  Port string
  AccessLog string
}

func startProxyServer(config *Config) {
  nginxDir := createDir(".hancho", "nginx")
  logsDir := createDir(nginxDir, "logs")
  serversDir := createDir(nginxDir, "servers")

  t := template.Must(template.New("nginx.conf").Parse(nginxTemplate))
  f := filepath.Join(nginxDir, "nginx.conf")
  a := abs(logsDir, "access.log")
  writeConf(f, t, NginxConf{a})

  t = template.Must(template.New("server.conf").Parse(serverTemplate))
  for name, service := range config.Services {
    host := name + ".hancho.localhost"
    a := abs(logsDir, host + ".access.log")
    f := filepath.Join(serversDir, host + ".conf")
    writeConf(f, t, ServerConf{host, service.Port, a})
  }

  sigs := make(chan os.Signal, 1)
  signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

  go func() {
    <-sigs
    fmt.Println()
  }()

  cmd := progressCommand("nginx", "-c", abs(nginxDir, "nginx.conf"))
  fmt.Println("Running services...")
  fmt.Println("Use Ctrl-C to stop")
  cmd.Wait()
  fmt.Println("Gracefully stopping services")
}

func writeConf(filename string, t *template.Template, data interface{}) {
  f, err := os.Create(filename)
  checkError(err)

  defer f.Close()

  wr := bufio.NewWriter(f)

  err = t.Execute(wr, data)
  checkError(err)

  wr.Flush()
}

func abs(elem ...string) string {
  path, err := filepath.Abs(filepath.Join(elem...))
  checkError(err)

  return path
}

func createDir(parent, name string) string {
  dir := filepath.Join(parent, name)
  err := os.MkdirAll(dir, os.ModePerm)
  checkError(err)

  return dir
}

const nginxTemplate = `
#user  nobody;
worker_processes  1;
daemon off;

#error_log  logs/error.log;
#error_log  logs/error.log  notice;
#error_log  logs/error.log  info;

#pid        logs/nginx.pid;


events {
    worker_connections  1024;
}


http {
    #include       mime.types;
    #default_type  application/octet-stream;

    log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"';

    access_log  {{.AccessLog}}  main;

    sendfile        on;
    #tcp_nopush     on;

    #keepalive_timeout  0;
    keepalive_timeout  65;

    #gzip  on;

    include servers/*;

    #server {
        #listen       8080;
        #server_name  localhost;

        #charset koi8-r;

        #access_log  logs/host.access.log  main;

        #location / {
        #    root   html;
        #    index  index.html index.htm;
        #}

        #error_page  404              /404.html;

        # redirect server error pages to the static page /50x.html
        #
        #error_page   500 502 503 504  /50x.html;
        #location = /50x.html {
        #    root   html;
        #}

        # proxy the PHP scripts to Apache listening on 127.0.0.1:80
        #
        #location ~ \.php$ {
        #    proxy_pass   http://127.0.0.1;
        #}

        # pass the PHP scripts to FastCGI server listening on 127.0.0.1:9000
        #
        #location ~ \.php$ {
        #    root           html;
        #    fastcgi_pass   127.0.0.1:9000;
        #    fastcgi_index  index.php;
        #    fastcgi_param  SCRIPT_FILENAME  /scripts$fastcgi_script_name;
        #    include        fastcgi_params;
        #}

        # deny access to .htaccess files, if Apache's document root
        # concurs with nginx's one
        #
        #location ~ /\.ht {
        #    deny  all;
        #}
    #}


    # another virtual host using mix of IP-, name-, and port-based configuration
    #
    #server {
    #    listen       8000;
    #    listen       somename:8080;
    #    server_name  somename  alias  another.alias;

    #    location / {
    #        root   html;
    #        index  index.html index.htm;
    #    }
    #}


    # HTTPS server
    #
    #server {
    #    listen       443 ssl;
    #    server_name  localhost;

    #    ssl_certificate      cert.pem;
    #    ssl_certificate_key  cert.key;

    #    ssl_session_cache    shared:SSL:1m;
    #    ssl_session_timeout  5m;

    #    ssl_ciphers  HIGH:!aNULL:!MD5;
    #    ssl_prefer_server_ciphers  on;

    #    location / {
    #        root   html;
    #        index  index.html index.htm;
    #    }
    #}
}
`

const serverTemplate = `
server {
    listen       8080;
    server_name  {{.Host}};

    access_log  {{.AccessLog}}  main;

    location / {
        proxy_pass   http://localhost:{{.Port}};
    }
}
`
