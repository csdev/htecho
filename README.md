# HTEcho

## Usage

HTEcho is a minimal HTTP echo server, useful for testing API clients and
network connectivity situations.

```bash
docker pull csang/htecho:0.1
docker run --rm -p 8080:8080 csang/htecho:0.1 --addr=0.0.0.0:8080
```

Use `--addr` to specify the server's bind address. Most systems also require publishing
the container port via `docker run -p` so the container is accessible from the host.

See the [Docker Hub repo](https://hub.docker.com/repository/docker/csang/htecho/general)
for available image tags and architectures.

## Behavior

### Basic Request Information

HTEcho responds to any HTTP request, regardless of HTTP method, URL path, or credentials.
The response status code is always `200 OK`.

The response body is equal to the request body. The response `Content-Type` is echoed from the request
header, or is autodetected if not present. `Content-Length` is set automatically.

Other information about the request is echoed back in the response headers:

* `X-Echo-Method`: The request method, e.g., `POST`.
* `X-Echo-Path`: The request URL path.
* `X-Echo-Query`: The raw query string (everything after the `?` in the URL).
* `X-Echo-Header-*`: Request headers are echoed back with this prefix, e.g., `X-Echo-Header-User-Agent`.

### Additional Request Information

To avoid leaking sensitive data, some information is excluded by default:

* Auth headers (use `--include-auth`)
    * `X-Echo-Header-Authorization`
    * `X-Echo-Header-Proxy-Authorization`
* IP address headers (use `--include-ips`)
    * `X-Echo-Header-X-Forwarded-For`
    * `X-Echo-Header-Forwarded`
    * `X-Echo-Addr` (the client's IP address, as seen by the server)

You can also use `-A` to include everything:

```bash
docker run --rm csang/htecho:latest -A
```

### Proxies

Deploying HTEcho behind a load balancer or proxy does not result in any special behavior.
In particular, there is no "proxy fix" middleware. `X-Echo-Addr` may return the IP address
of the proxy, and you may need to parse `X-Echo-Header-X-Forwarded-For` on your own.

As a network testing tool, HTEcho avoids parsing, validating, or transforming the request
it receives so you can focus on the raw data.

## Example

```bash
curl -v -u 'user:pass' -d '{"foo":"bar"}' -H 'Content-Type: application/json' \
    'http://localhost:8080/?foo=bar&bar=baz'
```
```
*   Trying 127.0.0.1:8080...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
* Server auth using Basic with user 'user'
> POST /?foo=bar&bar=baz HTTP/1.1
> Host: localhost:8080
> Authorization: Basic dXNlcjpwYXNz
> User-Agent: curl/7.68.0
> Accept: */*
> Content-Type: application/json
> Content-Length: 13
>
* upload completely sent off: 13 out of 13 bytes
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: application/json
< X-Echo-Header-Accept: */*
< X-Echo-Header-Content-Length: 13
< X-Echo-Header-Content-Type: application/json
< X-Echo-Header-User-Agent: curl/7.68.0
< X-Echo-Method: POST
< X-Echo-Path: /
< X-Echo-Query: foo=bar&bar=baz
< Date: Thu, 02 Feb 2023 22:26:38 GMT
< Content-Length: 13
<
* Connection #0 to host localhost left intact
{"foo":"bar"}
```

## Advanced Options

### Access Log

Log all requests to stdout:

```bash
docker run --rm csang/htecho:latest --access-log
```

Sample log format:

```
2023/02/06 16:57:51 htecho.server: listening on 0.0.0.0:8080
2023/02/06 16:58:27 htecho.request: 172.18.0.1:52422 GET /path/to/something?query (0 bytes)
2023/02/06 16:58:41 htecho.request: 172.18.0.1:52426 POST /path/to/something (13 bytes)
```

### Server Timeouts

HTEcho has 1 minute timeouts for reading the request and writing the response.
These values can be adjusted:

```bash
docker run --rm csang/htecho:latest --read-timeout=1m --write-timeout=1m
```

## Developer Workflow

Build and run the server locally using Docker Compose:

```bash
docker-compose run --build --rm --service-ports server
```

Available services:

* `server` - the production build
* `dev` - dev container with interactive shell for running arbitrary build and debug commands
* `test` - run tests in dev container
* `cover` - run tests with code coverage
