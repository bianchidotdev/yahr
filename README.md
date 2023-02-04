# Yahr

Yahr (YAml Http Request) is a tool for making HTTP requests based on YAML files.

It's goal is to provide a simple interface for testing APIs or retrieving data in a structure format (YAML)
that is also easy to share with friends and coworkers.

The primary motivation is to be able to build a set of requests for testing an API that can live
within the repo and be shared by a team of developers. Think Postman or Insomnia with git syncing
and a CLI interface.

> **Note**
> This is an alpha project and should not be considered stable

## Features

* Makes HTTP requests (GET, POST, PUT, DELETE)
* Supports dynamic path parameters
* Supports sending JSON data
* Supports custom headers
* Supports gotemplate in YAML files for more dynamic requests (ie. Auth tokens) + reading from `.env` file`

## Installation

``` sh
go install github.com/michaeldbianchi/yahr@latest
yahr -v
```


## Usage

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

yahr run -s httpbin get | jq .origin

# "89.188.181.42"

# with dynamic path parameters

yahr -c fixtures/github.yaml run -p owner=michaeldbianchi -p repo=yahr -s github get_repo | jq "{stargazers_count, open_issues}"

# {
#   "stargazers_count": 0,
#   "open_issues": 0
# }

```

### Scripting

Right now, `yahr` as a tool can only make a single request at a time.
However, if you are interested in stringing together requests, you can absolutely do that primitively
through using a shell script.

Check out [examples/series_of_requests.sh](examples/series_of_requests.sh) for a simple example.

In a future release, I expect to refine the `yahr` core library so it's
also possible to use this tool easily within go scripts.

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

  private_server:
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

* [x] Functional http requests based off yaml file - v0.1.0
* [ ] brew install
* [ ] Use of yahr as a go library (not just a cli app) - v0.2.0
  * [ ] with examples
* [ ] Environments for changing http config across a series of requests (dev, staging, prod) - v0.3.0
* [ ] Sequences of requests
* [ ] Inherit/import from other configuration files

Additional nice-to-haves:
* [ ] Output request as curl (maybe in yahr requests show GROUP NAME)

### Anti-features

This is a set of features I don't foresee ever implementing in this project, usually because I foresee the complexity overwhelming what should be a straightforward tool.

* Arbitrarily nestable groups for deduplication of auth/headers/base-urls
* Feeding output from one request into another (I think this would be too complex for the codebase and can be done easily with light scripting)

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

