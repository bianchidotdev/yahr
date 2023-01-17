# Yahr

Yahr (YAml Http Request) is a system for making HTTP requests based on YAML files.

> **Note**
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
    host: google.com
    scheme: https // default: https
    path: /
  private_server_get
    host: localhost
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


## Appendix

Potential future spec

``` yaml
requests:
  httpbin:
    host: httpbin.org
    requests:
      get:
        path: /get

      get200:
        path: /status/200

      post:
        path: /post
        # should we allow yaml map and translate it into json?
        # how do we deal with non json payloads
        payload: {"test_payload": "yep. this is a test"}

  httpbin2:
    requests:
      get:
        path: /get

environments:
  dev:
    scheme: http
    host: localhost:2222
    headers:
      Authorization: Bearer fake-api-key

  prod:
    headers:
      Authorization: Bearer {{ .SECRET_KEY_ENV_VAR }}
```

``` sh
yahr run httpbin get -e dev
# sends GET http://localhost:2222/get with fake-api-key bearer token

yahr run httpbin post -e prod
# sends POST https://httpbin.org/post with the sample payload and a bearer token from the env var SECRET_KEY_ENV_VAR

yahr run httpbin get200
# sends GET https://httpbin.org/status/200 with no bearer tokens

yahr run httpbin2 get -e dev
# sends GET http://localhost:2222/get with fake-api-key bearer token

yahr run httpbin2 get
# errors because it is missing a required host
```

