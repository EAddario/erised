
# erised
A simple **http server** to test arbitrary responses.

# Usage
Upon executing **erised** it will listen on port **8080** for incoming http requests.

HTTP verbs (e.g. GET, POST, PATCH, etc.), URL Paths and query strings are **ignored**. 

Response behaviour is controlled via custom http headers:

|Name|Purpose|
|--|--|
|X-Erised-Data|Returns the value **as is** in the response body|
|X-Erised-Content-Type|Returns the value **as is** in the *Content-Type* response header|
|X-Erised-Status-Code|Used to set the *http status code* value|
|X-Erised-Location|Returns the value **as is** of a new URL or path when 300 ≤ *X-Erised-Status-Code* < 310

By design, no validation is performed on *X-Erised-Data*, *X-Erised-Content-Type* or *X-Erised-Location*. Values are returned **as is**

Valid *X-Erised-Status-Code* values are as follows:
```text
OK or 200

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
< Content-Type: text/plain; charset=utf-8
< Date: Tue, 29 Dec 2020 18:35:48 GMT
< Content-Length: 0
<
* Connection #0 to host localhost left intact

* Closing connection 0
```

### Request returning *Hello World* in the response's body:
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
< Content-Type: text/plain; charset=utf-8
< Date: Tue, 29 Dec 2020 18:38:10 GMT
< Content-Length: 11
<
* Connection #0 to host localhost left intact
Hello World
* Closing connection 0
```

### Request returning *[{"Hello":"World"}]* in the response's body and *text/json* in the header's Content-Type
```
curl -w '\n' -v -H "X-Erised-Content-Type:text/json; charset=utf-8" -H "X-Erised-Data:[{"Hello":"World"}]" http://localhost:8080
```
```sh
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 8080 (#0)
> GET / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.64.1
> Accept: */*
> X-Erised-Content-Type:text/json; charset=utf-8
> X-Erised-Data:[{Hello:World}]
>
< HTTP/1.1 200 OK
< Content-Type: text/json; charset=utf-8
< Date: Tue, 29 Dec 2020 18:43:55 GMT
< Content-Length: 15
<
* Connection #0 to host localhost left intact
[{Hello:World}]
* Closing connection 0
```

### Request returning *text* in the response's body, an arbitrary *application/teapot* in the header's Content-Type and [*418 I'm a teapot*](https://save418.com/) in the header's Status Code
```
curl -w '\n' -v -H "X-Erised-Status-Code:Teapot" -H "X-Erised-Content-Type:application/teapot; charset=utf-8" -H "X-Erised-Data:Server refuses to brew coffee because it is, permanently, a teapot." http://localhost:8080
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
> X-Erised-Content-Type:application/teapot; charset=utf-8
> X-Erised-Data:Server refuses to brew coffee because it is, permanently, a teapot.
>
< HTTP/1.1 418 I'm a teapot
< Content-Type: application/teapot; charset=utf-8
< Date: Tue, 29 Dec 2020 18:54:46 GMT
< Content-Length: 67
<
* Connection #0 to host localhost left intact
Server refuses to brew coffee because it is, permanently, a teapot.
* Closing connection 0
```

# Release History
* WIP (not yet released) - Add HTTP redirection status codes (300's)
* v0.0.1 - Initial release

# Known Issues
**erised** is full of bugs and "_...men have wasted away before it, not knowing if what they have seen is real, or even possible..._" so use it with caution for it gives no knowledge or truth.

Of all its deficiencies, the most notable are:
* There are not tests (yet)
* **erised** offers no help
* Server parameters are hardcoded
* Server does not shutdown gracefully. To stop, process must be terminated
* https protocol is not supported
* **erised** does not scale well

I may or may not address any of this in a future release. Caveat Emptor

# Motivation
When developing and testing REST based API clients, sooner or later I'd come across situations where I needed a quick and easy way to dynamically test endpoint's responses under different scenarios. Although there are many excellent frameworks and mock servers available, the time and effort required to configure them is sometimes not justified, specially if the application under test provides 10's or 100's of paths, so after some brief and unsuccessful googling I decided to create my own.

**erised** was inspired somewhat by [Kenneth Reitz's](https://kennethreitz.org/) HTTP Request & Response Service [httpbin.io](https://httpbin.org/) and it may offer similar functionality in future releases.

The typical use case is to get a response to an arbitrary http request where the content of the body has a predetermined value and your ability to control the server's behaviour is limited or non-existent.

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

**Or**, you could use **erised** like this:
```sh
curl -w '\n' -v \
-H "X-Erised-Status-Code:OK" \
-H "X-Erised-Content-Type:application/json; charset=UTF-8" \
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
> X-Erised-Content-Type:application/json; charset=UTF-8
> X-Erised-Data:{"categories":[],"created_at":"2020-01-05 13:42:26.766831","icon_url":"https://assets.chucknorris.host/img/avatar/chuck-norris.png","id":"CfW0ccNFTpeq_v1r13IjTQ","updated_at":"2020-01-05 13:42:26.766831","url":"https://api.chucknorris.io/jokes/CfW0ccNFTpeq_v1r13IjTQ","value":"The lord giveth and Chuck Norris taketh away"}
>
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=UTF-8
< Date: Wed, 30 Dec 2020 01:13:54 GMT
< Content-Length: 323
<
* Connection #0 to host localhost left intact
{"categories":[],"created_at":"2020-01-05 13:42:26.766831","icon_url":"https://assets.chucknorris.host/img/avatar/chuck-norris.png","id":"CfW0ccNFTpeq_v1r13IjTQ","updated_at":"2020-01-05 13:42:26.766831","url":"https://api.chucknorris.io/jokes/CfW0ccNFTpeq_v1r13IjTQ","value":"The lord giveth and Chuck Norris taketh away"}
* Closing connection 0
```

**and** also to test some common failures like,
```sh
curl -w '\n' -v \
-H "X-Erised-Status-Code:NotFound" \
-H "X-Erised-Content-Type:application/json; charset=UTF-8" \
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
> X-Erised-Content-Type:application/json; charset=UTF-8
> X-Erised-Data:{"timestamp":"2020-12-30T11:21:32.793Z","status":404,"error":"Not Found","message":"Chuck Norris knows everything there is to know - Except where this page is.","path":"/jokes/random"}
>
< HTTP/1.1 404 Not Found
< Content-Type: application/json; charset=UTF-8
< Date: Wed, 30 Dec 2020 11:25:21 GMT
< Content-Length: 184
<
* Connection #0 to host localhost left intact
{"timestamp":"2020-12-30T11:21:32.793Z","status":404,"error":"Not Found","message":"Chuck Norris knows everything there is to know - Except where this page is.","path":"/jokes/random"}
* Closing connection 0
```
