# erised
A nimble **http server** to test arbitrary REST API responses.

# Usage
`erised [options]`
```text
Parameters:
  -idle int
    	maximum time in seconds to wait for the next request when keep-alive is enabled (default 120)
  -json
    	use JSON log format
  -level string
    	one of debug/info/warn/error/off (default "info")
  -path string
    	path to search recursively for X-Erised-Response-File (path is restricted to the directory or subdirectories where the program is invoked)
  -port int
    	port to listen (default 8080)
  -read int
    	maximum duration in seconds for reading the entire request (default 5)
  -write int
    	maximum duration in seconds before timing out response writes (default 10)
```

For help type **erised -h**

When executing **erised** with no parameters, the server will listen on port **8080** for incoming http requests.

If you're using the _-path_ option, please **EXERCISE GREAT CAUTION** when setting the path to search. See **Known Issues** for more information. For security reasons, path is restricted to the directory or subdirectories where the program was invoked.

The latest version is also available as a Docker image at [edaddario/erised](https://hub.docker.com/r/edaddario/erised).

To start the server in a docker container, with defaults values, execute the following command:

```sh
docker run --rm -p 8080:8080 --name erised edaddario/erised [options]
```

If you would like to return file based responses (_X-Erised-Response-File_ set) when using the docker image, you'll need to map the directory containing your local files and set the _-path_ option accordingly.

The following example maps the **/local_directory/response_files** directory in your local machine to **/files** in the docker image, and then sets the **-path** option:

```sh
docker run --rm -p 8080:8080 --name erised -v /local_directory/response_files:/files edaddario/erised -path ./files
```

URL routes, HTTP methods (e.g. GET, POST, PATCH, etc.), query strings and body are **ignored**, except for:

| Name            | Method | Purpose                           |
|-----------------|--------|-----------------------------------|
| erised/headers  | GET    | Returns request headers           |
| erised/info     | GET    | Returns miscellaneous information |
| erised/ip       | GET    | Returns the client IP             |
| erised/shutdown | POST   | Shutdowns the server              |

The `erised/echoserver` path will ignore any additional segments after `/echoserver`, including HTTP methods, query strings and body, and it will return a webpage displaying server information and the request's parameters.

| Name                | Method | Purpose                                                                      |
|---------------------|--------|------------------------------------------------------------------------------|
| erised/echoserver/* | any    | Returns a webpage displaying server information and the request's parameters |

Erised's response behaviour is controlled via custom headers in the http request:

| Name                    | Purpose                                                                                                                                                                                                                                                                                                              |
|-------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| X-Erised-Content-Type   | Sets the response _Content-Type_. Valid values are **text** (default) for _text/plain_, **json** for _application/json_, **xml** for _application/xml_ and **gzip** for _application/octet-stream_. When using **gzip**, _Content-Encoding_ is also set to **gzip** and the response body is compressed accordingly. |
| X-Erised-Data           | Returns the **same** value in the response body                                                                                                                                                                                                                                                                      |
| X-Erised-Headers        | Returns the value(s) in the response header. Values **must** be in a JSON key/value list                                                                                                                                                                                                                             |
| X-Erised-Location       | Sets the response _Location_ to the new (redirected) URL or path, when 300 ≤ _X-Erised-Status-Code_ < 310                                                                                                                                                                                                            |
| X-Erised-Response-Delay | Number of **milliseconds** to wait before sending response back to client                                                                                                                                                                                                                                            |
| X-Erised-Response-File  | Returns the contents of **file** in the response body. If present, _X-Erised-Data_ is ignored                                                                                                                                                                                                                        |
| X-Erised-Status-Code    | Sets the HTTP Status Code                                                                                                                                                                                                                                                                                            |

No validation is performed on _X-Erised-Data_ or _X-Erised-Location_.

Valid _X-Erised-Status-Code_ values are:
```text
OK or 200 (default)

MultipleChoices or 300
MovedPermanently or 301
Found or 302
SeeOther or 303
UseProxy or 305
TemporaryRedirect or 307
PermanentRedirect or 308

BadRequest or 400
Unauthorized or 401
PaymentRequired or 402
Forbidden or 403
NotFound or 404
MethodNotAllowed or 405
RequestTimeout or 408
Conflict or 409
Gone or 410
Teapot or 418
TooManyRequests or 429

InternalServerError or 500
NotImplemented or 501
BadGateway or 502
ServiceUnavailable or 503
GatewayTimeout or 504
HTTPVersionNotSupported or 505
InsufficientStorage or 507
LoopDetected or 508
NotExtended or 510
NetworkAuthenticationRequired or 511
```
Any other value will resolve to 200 (OK)

# Release History
* v0.9.6 - Rename _erised/webpage_ to _erised/echoserver_ and add headers and server environment information
* v0.8.3 - Add _erised/webpage_ path, add multi-architecture docker images, minor refactoring, and minor cosmetic changes
* v0.7.0 - Improve response file processing and security, change logging type, and minor source code readability changes
* v0.6.11 - Further server shutdown improvements, minor efficiency improvements, general code refactoring and bug fixes
* v0.6.7 - Improve server shutdown handling, and restrict allowed methods for _erised/headers_, _erised/ip_, _erised/info_ and _erised/shutdown_ routes  
* v0.5.4 - Update dependencies
* v0.5.3 - Add file based responses
* v0.4.1 - Add route concurrency, update tests and dependencies
* v0.3.4 - Add [gomega](https://onsi.github.io/gomega/) assertion library, refactor tests to use Ω assertions and minor bug fixes
* v0.3.0 - Add [goblin](https://github.com/franela/goblin) framework and unit tests
* v0.2.5 - Switch to zerolog logging framework, add erised/shutdown path
* v0.2.2 - Add custom headers, add dockerfile
* v0.2.1 - Add gzip compression, improve erised/headers json handling
* v0.0.3 - Add erised/headers, erised/ip and erised/info routes. Add delayed responses
* v0.0.2 - Add HTTP redirection status codes (300's), startup configuration parameters and request's logging
* v0.0.1 - Initial release

# Known Issues
**erised** may be full of bugs. Poeple "_... have wasted away before it, not knowing if what they have seen is real, or even possible..._" so, use it with caution for it gives no knowledge or truth.

Of all of its deficiencies, the most notable is:
* Using the _-path_ option could lead to significant security risks. When the _X-Erised-Response-File_ header is set it will search recursively for a matching filename in the current directory or **all** subdirectories underneath, returning the contents of the first match. For security reasons, path is restricted to the directory or subdirectories where the program was invoked.
* https protocol is not yet supported

I may or may not address these issues in a future release. Caveat Emptor

# Motivation
When developing and testing REST API clients, sooner or later I'd come across situations where I needed a quick and easy way to dynamically test endpoint's responses under different scenarios. Although there are many excellent frameworks and mock servers available, the time and effort required to configure them is sometimes not justified, specially if the application under test exposes many routes, so after some brief and unsuccessful googling I decided to create my own.

**erised** was inspired by [Kenneth Reitz's](https://kennethreitz.org/) HTTP Request & Response Service [httpbin.io](https://httpbin.org/) and it may offer similar functionality in future releases.

The typical use case is to get a response to an arbitrary http request when your ability to control the server's behaviour is limited or non-existent.

Imagine you're developing some client for [api.chucknorris.io](https://api.chucknorris.io/) and want to test the **/jokes/random** path. You could certainly make live calls against the server:
```sh
curl -w '\n' -v -k https://api.chucknorris.io/jokes/random
```
(response edited for clarity)
```sh
*   Trying 104.31.94.71...
* TCP_NODELAY set
* Connected to api.chucknorris.io (104.31.94.71) port 443 (#0)
> GET /jokes/random HTTP/2
> Host: api.chucknorris.io
> User-Agent: curl/7.64.1
> Accept: */*
>
< HTTP/2 200
< date: Wed, 30 Dec 2020 00:21:14 GMT
< content-type: application/json;charset=UTF-8
<
* Connection #0 to host api.chucknorris.io left intact
{"categories":[],"created_at":"2020-01-05 13:42:18.823766","icon_url":"https://assets.chucknorris.host/img/avatar/chuck-norris.png","id":"CfW0ccNFTpeq_v1r13IjTQ","updated_at":"2020-01-05 13:42:18.823766","url":"https://api.chucknorris.io/jokes/CfW0ccNFTpeq_v1r13IjTQ","value":"The lord giveth and Chuck Norris taketh away"}
* Closing connection 0
```

**Or**, even better yet, you could use **erised** like this:
```sh
curl -w '\n' -v \
-H "X-Erised-Status-Code:OK" \
-H "X-Erised-Content-Type:json" \
-H "X-Erised-Data:{\"categories\":[],\"created_at\":\"2020-01-05 13:42:26.766831\",\"icon_url\":\"https://assets.chucknorris.host/img/avatar/chuck-norris.png\",\"id\":\"CfW0ccNFTpeq_v1r13IjTQ\",\"updated_at\":\"2020-01-05 13:42:26.766831\",\"url\":\"https://api.chucknorris.io/jokes/CfW0ccNFTpeq_v1r13IjTQ\",\"value\":\"The lord giveth and Chuck Norris taketh away\"}" \
http://localhost:8080/jokes/random
```
```sh
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 8080 (#0)
> GET /jokes/random HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.64.1
> Accept: */*
> X-Erised-Status-Code:OK
> X-Erised-Content-Type:json
> X-Erised-Data:{"categories":[],"created_at":"2020-01-05 13:42:26.766831","icon_url":"https://assets.chucknorris.host/img/avatar/chuck-norris.png","id":"CfW0ccNFTpeq_v1r13IjTQ","updated_at":"2020-01-05 13:42:26.766831","url":"https://api.chucknorris.io/jokes/CfW0ccNFTpeq_v1r13IjTQ","value":"The lord giveth and Chuck Norris taketh away"}
>
< HTTP/1.1 200 OK
< Content-Encoding: identity
< Content-Type: application/json
< Date: Wed, 30 Dec 2020 01:13:54 GMT
< Content-Length: 323
<
* Connection #0 to host localhost left intact
{"categories":[],"created_at":"2020-01-05 13:42:26.766831","icon_url":"https://assets.chucknorris.host/img/avatar/chuck-norris.png","id":"CfW0ccNFTpeq_v1r13IjTQ","updated_at":"2020-01-05 13:42:26.766831","url":"https://api.chucknorris.io/jokes/CfW0ccNFTpeq_v1r13IjTQ","value":"The lord giveth and Chuck Norris taketh away"}
* Closing connection 0
```

**and** even simulate common failures like,
```sh
curl -w '\n' -v \
-H "X-Erised-Status-Code:NotFound" \
-H "X-Erised-Content-Type:json" \
-H "X-Erised-Data:{\"timestamp\":\"2020-12-30T11:21:32.793Z\",\"status\":404,\"error\":\"Not Found\",\"message\":\"Chuck Norris knows everything there is to know - Except where this page is.\",\"path\":\"/jokes/random\"}" \
http://localhost:8080/jokes/random
```
```sh
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 8080 (#0)
> GET /jokes/random HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.64.1
> Accept: */*
> X-Erised-Status-Code:NotFound
> X-Erised-Content-Type:json
> X-Erised-Data:{"timestamp":"2020-12-30T11:21:32.793Z","status":404,"error":"Not Found","message":"Chuck Norris knows everything there is to know - Except where this page is.","path":"/jokes/random"}
>
< HTTP/1.1 404 Not Found
< Content-Encoding: identity
< Content-Type: application/json
< Date: Wed, 30 Dec 2020 11:25:21 GMT
< Content-Length: 184
<
* Connection #0 to host localhost left intact
{"timestamp":"2020-12-30T11:21:32.793Z","status":404,"error":"Not Found","message":"Chuck Norris knows everything there is to know - Except where this page is.","path":"/jokes/random"}
* Closing connection 0
```

# Examples
### Simple request returning nothing in the response's body:
```
curl -w '\n' -v http://localhost:8080
```
```sh
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 8080 (#0)
> GET / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.64.1
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Encoding: identity
< Content-Type: text/plain
< Date: Tue, 29 Dec 2020 18:35:48 GMT
< Content-Length: 0
<
* Connection #0 to host localhost left intact

* Closing connection 0
```

### Simple request returning custom headers only:
```
curl -w '\n' -I -H "X-Erised-Headers:{\"My-Header\":\"Hello World\",\"Another-Header\":\"Goodbye World\"}" http://localhost:8080
```
```sh
HTTP/1.1 200 OK
Another-Header: Goodbye World
Content-Encoding: identity
Content-Type: text/plain
My-Header: Hello World
Date: Sat, 13 Mar 2021 22:56:09 GMT
```
### Request returning _Hello World_ in the response's body:
```
curl -w '\n' -v -H "X-Erised-Data:Hello World" http://localhost:8080
```
```sh
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 8080 (#0)
> GET / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.64.1
> Accept: */*
> X-Erised-Data:Hello World
>
< HTTP/1.1 200 OK
< Content-Encoding: identity
< Content-Type: text/plain
< Date: Tue, 29 Dec 2020 18:38:10 GMT
< Content-Length: 11
<
* Connection #0 to host localhost left intact
Hello World
* Closing connection 0
```

### Request returning _[{"Hello":"World"}]_ in the response's body and _json_ in the header's Content-Type
```
curl -w '\n' -v -H "X-Erised-Content-Type:json" -H "X-Erised-Data:[{\"Hello\":\"World\"}]" http://localhost:8080
```
```sh
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 8080 (#0)
> GET / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.64.1
> Accept: */*
> X-Erised-Content-Type:json
> X-Erised-Data:[{Hello:World}]
>
< HTTP/1.1 200 OK
< Content-Encoding: identity
< Content-Type: application/json
< Date: Tue, 29 Dec 2020 18:43:55 GMT
< Content-Length: 15
<
* Connection #0 to host localhost left intact
[{"Hello":"World"}]
* Closing connection 0
```

### Request returning _text_ in the response body and [_418 I'm a teapot_](https://save418.com/) in the header's Status Code
```
curl -w '\n' -v -H "X-Erised-Status-Code:Teapot" -H "X-Erised-Data:Server refuses to brew coffee because it is, permanently, a teapot." http://localhost:8080
```
```sh
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 8080 (#0)
> GET / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.64.1
> Accept: */*
> X-Erised-Status-Code:Teapot
> X-Erised-Data:Server refuses to brew coffee because it is, permanently, a teapot.
>
< HTTP/1.1 418 I'm a teapot
< Content-Encoding: identity
< Content-Type: text/plain
< Date: Tue, 29 Dec 2020 18:54:46 GMT
< Content-Length: 67
<
* Connection #0 to host localhost left intact
Server refuses to brew coffee because it is, permanently, a teapot.
* Closing connection 0
```
