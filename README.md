# Yahr

Yahr (YAml Http Request) is a system for making HTTP requests based on YAML files.

> **NOTE**
> This is an alpha project and should not be considered stable

## Getting started

Install with `go install`

``` sh
go install github.com/michaeldbianchi/yahr
yahr version
```

## Reference

YAML spec

``` yaml
requests:
  google_get:
    url: google.com
    scheme: https // default: https
    path: /
  private_server_get
    url: localhost
    port: 2222
    scheme: http
    path: /opl/health
    headers:
      Authorization: Bearer {{ .PRIVATE_SERVER_ENV_VAR }}_
```

## Roadmap

* [ ] Functional http requests based off yaml file
* [ ] Arbitrarily nestable groups for deduplication of auth/headers/base-urls
* [ ] Inherit/import from other configuration files


