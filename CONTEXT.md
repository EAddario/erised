# Erised

A nimble HTTP and echo server for testing arbitrary HTTP requests and REST API responses. Clients control response behaviour entirely through custom request headers — no configuration files or endpoint setup required.

## Language

**Erised**:
An HTTP/HTTPS test server that returns responses dictated by `X-Erised-*` request headers.
_Avoid_: Mock server, stub server

**Erised Response Header**:
Any request header prefixed with `X-Erised-` that controls response behaviour (status code, content type, body, delay, redirects, etc.).

**Erised Data**:
The `X-Erised-Data` header value, returned verbatim as the response body. Ignored when `X-Erised-Response-File` is present.

**Erised Response File**:
The `X-Erised-Response-File` header value — a filename searched recursively under `-path` and returned as the response body. First match wins.

**Erised Status Code**:
The `X-Erised-Status-Code` header value. Accepts either numeric codes (`200`, `404`, `500`) or human-readable names (`OK`, `NotFound`, `Teapot`). Falls back to `200` on unrecognized values.

**Erised Content Type**:
The `X-Erised-Content-Type` header value: `text` (text/plain), `json` (application/json), `xml` (application/xml), `gzip` (application/octet-stream + gzip encoding), or `html` (text/html).

**Erised Headers**:
The `X-Erised-Headers` header value — a JSON key/value object injected as response headers.

**Erised Location**:
The `X-Erised-Location` header value, set as the response `Location` header when `300` ≤ status < `310`.

**Erised Response Delay**:
The `X-Erised-Response-Delay` header value in milliseconds, instructing the server to sleep before responding.

**Erised Intent**:
The structured, purely logical representation of the desired HTTP response, derived from parsing the **Erised Response Headers**.

**Echo Server**:
The `/erised/echoserver/*` route returning an HTML page displaying server environment variables, request headers, method, path, body, and timing.

**Landing Handler**:
The catch-all `/` handler processing all `X-Erised-*` headers for arbitrary paths, methods, and bodies.

## Relationships

- Every request to **Landing Handler** may carry any combination of **Erised Response Headers**
- An **Erised Response File** supersedes **Erised Data** when both are present
- An **Erised Location** only takes effect when **Erised Status Code** is a redirect (30x)
- An **Erised Response Delay** pauses all response processing for the specified duration
- The **Echo Server** ignores `X-Erised-*` headers entirely — it always returns its own HTML page

## Built-in routes

| Route                     | Method | Purpose                                       |
| ------------------------- | ------ | --------------------------------------------- |
| `/`                       | any    | Landing Handler — processes `X-Erised-*` headers |
| `/erised/headers`         | GET    | Returns all request headers as JSON           |
| `/erised/info`            | GET    | Returns host, method, protocol, and URI       |
| `/erised/ip`              | GET    | Returns client remote address                 |
| `/erised/shutdown`        | POST   | Gracefully shuts down the server              |
| `/erised/echoserver/*`    | any    | Echo Server — diagnostic HTML page            |

## Example dialogue

> **Dev:** "I want to test a 404 JSON response with a delay. Do I need to configure routes?"
> **Domain expert:** "No — every request hits the **Landing Handler**. Just set `X-Erised-Status-Code:NotFound`, `X-Erised-Content-Type:json`, and `X-Erised-Response-Delay:500` on any request. No route setup needed."

## Flagged ambiguities

- "profile" / "trace" — separate flags for CPU profiling and execution tracing; both write to files with `.prf` and `.trc` extensions respectively, but neither has a validation step if the file can't be created (falls back to a warning).