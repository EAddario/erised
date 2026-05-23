# erised
A nimble **HTTP** and **echo** server designed to test arbitrary HTTP requests and mock REST API responses.

# Usage
`erised [options]`

```text
Parameters:
  -cert string
    	path to a valid X.509 certificate file
  -https
    	use HTTPS instead of HTTP. A valid X.509 certificate and private key are required
  -idle int
    	maximum time in seconds to wait for the next request when keep-alive is enabled (default 120)
  -json
    	use JSON log format
  -key string
    	path to a valid private key file
  -level string
    	one of debug/info/warn/error/off (default "info")
  -path string
    	path to search recursively for X-Erised-Response-File
  -port int
    	port to listen. Default is 8080 for HTTP and 8443 for HTTPS
  -profile string
    	profile this session. A valid file name is required
  -read int
    	maximum duration in seconds for reading the entire request (default 5)
  -write int
    	maximum duration in seconds before timing out response writes (default 10)
```

*For help guidelines, type `erised -h`*.

By default, executing `erised` without any parameters starts the server listening on port `8080` for incoming HTTP requests.

When using the `-path` option, **EXERCISE GREAT CAUTION** regarding the target directory you select. For security reasons, file searches are strictly restricted to the directory (and its subdirectories) from which the program was invoked. Please see the **Known Issues** section below for more details.

The latest version of the server is also available as a Docker image at [edaddario/erised](https://hub.docker.com/r/edaddario/erised).

To spin up the server inside a Docker container using default values, run the following command:

```sh
docker run --rm -p 8080:8080 --name erised edaddario/erised
```

If you want the Docker container to serve file-based responses via the `X-Erised-Response-File` header, you must mount your local directory containing those files and configure the `-path` option accordingly.

The example below maps your local machine's `/local_directory/response_files` directory to `/files` inside the Docker image, then specifies that folder using the `-path` option:

```sh
docker run --rm -p 8080:8080 --name erised -v /local_directory/response_files:/files edaddario/erised -path ./files
```

By default, all URL routes, HTTP methods (e.g., GET, POST, PATCH), query strings, and request bodies are **ignored**, with the exception of the following built-in endpoints:

| Name            | Method | Purpose                                  |
|-----------------|--------|------------------------------------------|
| erised/headers  | GET    | Returns request headers                  |
| erised/info     | GET    | Returns miscellaneous server information |
| erised/ip       | GET    | Returns the client's IP address          |
| erised/shutdown | POST   | Gracefully shuts down the server         |

The `erised/echoserver` path acts as a wildcard catcher. It ignores any additional segments appended after `/echoserver` (including HTTP methods, query strings, and payloads) and returns an HTML page displaying environment information alongside the incoming request's parameters.

| Name                | Method | Purpose                                                                      |
|---------------------|--------|------------------------------------------------------------------------------|
| erised/echoserver/* | any    | Returns a webpage displaying server information and the request's parameters |

You can dynamically control Erised's response behavior by supplying custom headers within your HTTP request:

| Name                    | Purpose                                                                                                                                                                                                                                                                                                                                      |
|-------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| X-Erised-Content-Type   | Configures the response `Content-Type`. Valid identifiers are **text** (default) for `text/plain`, **json** for `application/json`, **xml** for `application/xml`, and **gzip** for `application/octet-stream`. When utilizing **gzip**, the response body is automatically compressed and the `Content-Encoding` header is set to **gzip**. |
| X-Erised-Data           | Returns the exact string value provided in the response body.                                                                                                                                                                                                                                                                                |
| X-Erised-Headers        | Injects custom header values into the response. Values **must** be provided as a JSON object of key/value pairs.                                                                                                                                                                                                                             |
| X-Erised-Location       | Configures the response `Location` header to redirect to a new URL or path, applicable when 300 ≤ X-Erised-Status-Code < 310.                                                                                                                                                                                                                |
| X-Erised-Response-Delay | Specifies the number of **milliseconds** the server should wait before returning the response to the client.                                                                                                                                                                                                                                 |
| X-Erised-Response-File  | Instructs the server to return the raw contents of the designated **file** in the response body. When this header is present, `X-Erised-Data` is ignored.                                                                                                                                                                                    |
| X-Erised-Status-Code    | Sets the HTTP Status Code returned by the server.                                                                                                                                                                                                                                                                                            |

*Note: No format validation is performed on the values passed to `X-Erised-Data` or `X-Erised-Location`*.

The supported values for `X-Erised-Status-Code` are:

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

Any unsupported or unrecognized status value will automatically default to `200 OK`.

# Release History
* **v0.23.3** - Refactored codebase and updated dependencies.
* **v0.11.2** - Added HTTPS capabilities, included built-in test certificates, added program execution metrics, and introduced a profiling option. Refactored internal variable names for better readability and replaced application panics with user-friendly fatal logs.
* **v0.9.7** - Refactored error handling architecture.
* **v0.9.6** - Renamed the `erised/webpage` route to `erised/echoserver` and appended headers along with server environment data to the output.
* **v0.8.3** - Added the initial `erised/webpage` route, implemented multi-architecture Docker image builds, and introduced minor refactoring and cosmetic adjustments.
* **v0.7.0** - Optimized file-based response handling and security, transitioned to a new logging format, and completed minor source code readability adjustments.
* **v0.6.11** - Implemented further optimizations for server shutdown sequences, introduced minor efficiency improvements, and completed general code refactoring and bug fixes.
* **v0.6.7** - Enhanced server shutdown handling and restricted permissible HTTP methods for the `erised/headers`, `erised/ip`, `erised/info`, and `erised/shutdown` endpoints.
* **v0.5.4** - Updated third-party dependencies.
* **v0.5.3** - Added capability for file-based HTTP responses.
* **v0.4.1** - Introduced route concurrency capabilities alongside test and dependency updates.
* **v0.3.4** - Integrated the [gomega](https://onsi.github.io/gomega/) assertion library, refactored existing unit tests to leverage Ω assertions, and resolved minor bugs.
* **v0.3.0** - Integrated the [goblin](https://github.com/franela/goblin) test framework and introduced core unit tests.
* **v0.2.5** - Switched the logging engine to the zerolog framework and added the `erised/shutdown` endpoint.
* **v0.2.2** - Introduced support for response customization via custom headers and added the initial Dockerfile.
* **v0.2.1** - Added native gzip compression capabilities and optimized JSON handling for the `erised/headers` route.
* **v0.0.3** - Added the `erised/headers`, `erised/ip`, and `erised/info` routes, and introduced support for delayed responses.
* **v0.0.2** - Added HTTP redirection status codes (300 series), configuration parameters for server startup, and basic request logging.
* **v0.0.1** - Initial project release.

# Enabling HTTPS
To activate HTTPS mode, `erised` requires a valid X.509 certificate signed by a globally trusted Certificate Authority (CA), such as IdenTrust, DigiCert, or Let's Encrypt.

If you do not possess an official certificate or prefer to validate connections locally, a root CA file (`localCA.pem`) is bundled inside the [/certs](https://www.google.com/search?q=./cmd/erised/certs) directory. Installing this file allows your operating system to trust the provided local testing certificate and key pair. While walking through the installation process for CA certificates on various platforms is outside the scope of this README, detailed guides can be easily retrieved via standard search engines or AI assistants. Once installed, verify that it appears as **Erised Test CA** and is explicitly marked as **trusted**.

Once your certificates are configured, launch `erised` in HTTPS mode by running:
`erised -https -cert erised.crt -key erised.key`

In this command, `erised.crt` represents your local site's X.509 certificate, and `erised.key` represents its matching private key.

### A word of caution about trusting certificates with unclear provenance:

While defining the mechanics of establishing cryptographically secure digital identities and archiving key-generation procedures is beyond the scope of this file, it is vital to emphasize the security implications of trusting unauthorized digital certificates. Beyond establishing secure, encrypted TLS communication channels, trusted root certificates can potentially be utilized to sign executable code or applications, allowing them to run with elevated administrative privileges on your system.

The Secure Sockets Layer (SSL) certificates packaged with `erised` are bound strictly to `localhost` and are explicitly intended to facilitate safe, encrypted loops within your own computer.

For complete transparency, the high-level workflow utilized to construct these assets is outlined below:

1. A local private key (`localCA.key`) was created to sign the Root CA certificate.
2. A Root CA certificate (`localCA.pem`) was generated and signed using the aforementioned CA key.
3. A distinct "site" private key (`erised.key`) was established.
4. An intermediate Certificate Signing Request (`erised.csr`) was produced and signed using the site's private key.
5. An X.509 V3 extension configuration file (`erised.ext`) was constructed to safely link the final certificate to the `localhost` domain.
6. The final site certificate (`erised.crt`) was compiled using the Root CA certificate, the CA private key, the intermediate CSR, and the V3 configuration file.

Please note that neither of the included private keys is protected by a passphrase. While this is a highly discouraged practice in production environments, it was done intentionally here to simplify local configuration and encourage tinkering.

# Known Issues

`erised` may contain bugs. True to its mythological namesake, *people have wasted away before it, not knowing if what they have seen is real, or even possible*. Use it with care, as it gives neither knowledge nor truth.

Its most prominent operational limitation is:

* Utilizing the `-path` option presents significant security risks. When the `X-Erised-Response-File` header is invoked, the engine recursively crawls the base directory and **all** underlying subdirectories to find a filename match, delivering the contents of the first asset it hits. To protect the host environment, file searching is strictly containerized to the immediate directory or subdirectories from which the server binary was originally executed.

Resolutions for this or any other pending anomalies may or may not be introduced in future updates. Refer to the [Caveat Emptor](https://www.google.com/search?q=./LICENSE) terms for details.

# Motivation

When building and debugging REST API clients, developers inevitably face scenarios requiring a fast, lightweight mechanism to test client behavior against dynamic server responses. While excellent enterprise mock servers exist, configuring them can demand time and effort that is difficult to justify—particularly when an application exposes a vast array of routes. After a brief and unsuccessful search for an out-of-the-box solution, this utility was born.

`erised` was partially inspired by Kenneth Reitz's HTTP Request & Response Service ([httpbin.io](https://httpbin.org/)) and subsequently by Marchandise Rudy's [Echo-Server](https://github.com/Ealenn/Echo-Server).

The classic use case involves simulating explicit backend payloads or errors for an arbitrary HTTP request when your capacity to control the live target server is limited or nonexistent.

For example, imagine you are developing a client wrapper for [api.chucknorris.io](https://api.chucknorris.io/) and want to validate your app's handling of the `/jokes/random` endpoint. You could make direct live calls against the real server:
```sh
curl -w '\n' -v -k https://api.chucknorris.io/jokes/random
```

*(Response truncated for clarity)*

```sh
* Trying 104.31.94.71...
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

**Alternatively**, you can use `erised` to mirror that successful response locally without hitting the external network:
```sh
curl -w '\n' -v \
-H "X-Erised-Status-Code:OK" \
-H "X-Erised-Content-Type:json" \
-H "X-Erised-Data:{\"categories\":[],\"created_at\":\"2020-01-05 13:42:26.766831\",\"icon_url\":\"https://assets.chucknorris.host/img/avatar/chuck-norris.png\",\"id\":\"CfW0ccNFTpeq_v1r13IjTQ\",\"updated_at\":\"2020-01-05 13:42:26.766831\",\"url\":\"https://api.chucknorris.io/jokes/CfW0ccNFTpeq_v1r13IjTQ\",\"value\":\"The lord giveth and Chuck Norris taketh away\"}" \
http://localhost:8080/jokes/random
```

```sh
* Trying ::1...
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

**And** you can seamlessly simulate common downstream failures, such as a missing resource error, to verify your application's resiliency:
```sh
curl -w '\n' -v \
-H "X-Erised-Status-Code:NotFound" \
-H "X-Erised-Content-Type:json" \
-H "X-Erised-Data:{\"timestamp\":\"2020-12-30T11:21:32.793Z\",\"status\":404,\"error\":\"Not Found\",\"message\":\"Chuck Norris knows everything there is to know - Except where this page is.\",\"path\":\"/jokes/random\"}" \
http://localhost:8080/jokes/random
```

```sh
* Trying ::1...
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
```sh
curl -w '\n' -v http://localhost:8080
```

```sh
* Trying ::1...
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
```sh
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

### Request returning *Hello World* in the response's body:
```sh
curl -w '\n' -v -H "X-Erised-Data:Hello World" http://localhost:8080
```

```sh
* Trying ::1...
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

### Request returning *[{"Hello":"World"}]* in the response's body and *json* in the header's Content-Type
```sh
curl -w '\n' -v -H "X-Erised-Content-Type:json" -H "X-Erised-Data:[{\"Hello\":\"World\"}]" http://localhost:8080
```

```sh
* Trying ::1...
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

### Request returning *text* in the response body and *[418 I'm a teapot](https://save418.com/)* in the header's Status Code
```sh
curl -w '\n' -v -H "X-Erised-Status-Code:Teapot" -H "X-Erised-Data:Server refuses to brew coffee because it is, permanently, a teapot." http://localhost:8080
```

```sh
* Trying ::1...
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
