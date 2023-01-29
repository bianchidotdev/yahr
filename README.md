# Yahr

Yahr (YAml Http Request) is a system for making HTTP requests based on YAML files.

> **Note**
> This is an alpha project and should not be considered stable

## MVP TODO list

* [x] Implement request list command
* [x] Implement run command to actually run requests
* [x] Fix config flag
* [x] Fix silent flag
* [x] More sophisticated versioning
* [x] Separate out business logic from CLI commands
* [x] Rework yaml spec to support groups of requests
* [x] Fix/implement non-GET methods
* [x] Rip out cobra
* [x] Env var templating
* [x] Make request methods more strict
* [ ] Implement at least happy-path tests
* [ ] Dynamic path variables `/users/:user_id`

## Getting started

Install with `go install`

``` sh
go install github.com/michaeldbianchi/yahr
yahr version
```

List requests and run them:

``` sh
yahr requests list

# +---------+--------+--------+--------------------------------+
# | Group   | Name   | Method | Endpoint                       |
# +---------+--------+--------+--------------------------------+
# | httpbin | get200 | get    | https://httpbin.org/status/200 |
# | httpbin | post   | post   | https://httpbin.org/post       |
# | httpbin | get    | get    | https://httpbin.org/get        |
# | local   | get    | get    | http://localhost:8080/get      |
# +---------+--------+--------+--------------------------------+

yahr run httpbin get

# Request:
# GET /get HTTP/1.1
# Host: httpbin.org
# User-Agent: Go-http-client/1.1
# Accept-Encoding: gzip
# 
# Status: 200
# Response Body:
#  {
#   "args": {},
#   "headers": {
#     "Accept-Encoding": "gzip",
#     "Host": "httpbin.org",
#     "User-Agent": "Go-http-client/2.0",
#     "X-Amzn-Trace-Id": "Root=1-63cb33d6-0e6501db48f84bf2260e7dc3"
#   },
#   "origin": "89.187.180.41",
#   "url": "https://httpbin.org/get"
# }

# use the -s silent flag for piping output

yahr run httpbin get -s | jq .origin

# "89.188.181.42"
```

## Reference

YAML spec

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

  private_server
    host: localhost
    port: 2222
    scheme: http
    headers:
      Authorization: Bearer {{ .PRIVATE_SERVER_ENV_VAR }}_
    requests:
      get:
        path: /opl/health
```

## Roadmap

* [ ] Functional http requests based off yaml file
* [ ] Use of yahr as a go library (not just a cli app)
* [ ] Sequences of requests
* [ ] Inherit/import from other configuration files

### Anti-features

This is a set of features I don't foresee ever implementing in this project, usually because I foresee the complexity overwhelming what should be a straightforward tool.

* [ ] Arbitrarily nestable groups for deduplication of auth/headers/base-urls
* [ ] Feeding output from one request into another (I think this would be too complex for the codebase and can be done easily with light scripting)

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

# TODO: spec out sequences
sequences:

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

