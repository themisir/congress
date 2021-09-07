# Congress

Ingress like routing and DNS service for microservices development.

## Installation

// TODO: Add installation instructions

# Configuration

Create a `congress.yaml` file on your project root directory. The file will
contain all the ingress like routes for our microservices. The file syntax
is similar to kubernetes ingress scheme.

```yaml
congress:
  ip: 127.0.0.1
  proxy:
    enabled: true
    port: 80
  dns:
    enabled: true
    port: 23
    fallback: 8.8.8.8
rules:
  - host: myapp.dev
    defaultBackend: http://frontend-service:80/
  - host: api.myapp.dev
    paths:
      - path: /catalog
        backend: http://catalog-service:80/
      - path: /checkout
        backend: http://checkout-service:80/
```

### `congress.ip`

IP address of the environment running congress instance. The IP address used to
respond DNS queries.

### `congress.proxy.enabled`

_(Default: true)_

Sets whether or not congress reverse-proxy should be enabled for routing
requests.

### `congress.proxy.port`

Port number for congress reverse-proxy server.

### `congress.dns.enabled`

_(Default: true)_

Sets whether or not congress DNS should be enabled. The DNS will respond to
questions for `rules.host` with `congress.ip`. The DNS could be used as default
name resolver on development environments.

### `congress.dns.port`

Port number for congress DNS.

### `congress.dns.fallback`

Fallback DNS server address. The queries will be sent to fallback DNS when
there's no configured rule available for the queried hostname.

