# HTEcho

## Usage

HTEcho is a minimal HTTP echo server, useful for testing API clients and
network connectivity situations.

Start the server on `localhost:8080`:

```bash
docker run --rm csang/htecho:latest
```

Or with a custom bind address and port:

```bash
docker run --rm csang/htecho:latest --addr=0.0.0.0:8000
```

## Behavior

HTEcho responds to any HTTP request, regardless of HTTP method, URL path, or credentials.
The response status code is always `200 OK`.

The response body is equal to the request body. The response `Content-Type` is echoed from the request
header, or is autodetected if not present. `Content-Length` is set automatically.

Request headers are echoed in the response, prefixed with `X-Echo-Header-`. To avoid unintentionally
leaking credentials, the `Authorization` and `Proxy-Authorization` headers are omitted. To echo them,
start the server with the `--include-auth` option.

Additionally, the request method and query string are available as `X-Echo-Method` and `X-Echo-Query`,
respectively.

Example:

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
< X-Echo-Query: foo=bar&bar=baz
< Date: Thu, 02 Feb 2023 22:26:38 GMT
< Content-Length: 13
<
* Connection #0 to host localhost left intact
{"foo":"bar"}
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
