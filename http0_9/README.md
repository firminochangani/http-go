# HTTP 0.9

This directory contains the implementation of an HTTP 0.9 client and server implemented from scratch based on the specs of HTTP version 0.9.
Such version is no longer supported by any modern HTTP client or browser therefore I had to write a bespoke [client](./client) to interact with the server.

## Specs

- [x] The client makes a TCP-IP connection to the host using the domain name or IP number , and the port number given in the address.
- [x] If the port number is not specified, 80 is always assumed for HTTP.
- [x] This request consists of the word "GET", a space, the document address, omitting the "http:, host and port parts when they are the coordinates just used to make the connection. (If a gateway is being used, then a full document address may be given specifying a different naming scheme).
- [x] The document address will consist of a single word (ie no spaces). If any further words are found on the request line, they MUST either be ignored, or else treated according to the full HTTP spec.
- [x] The response to a simple GET request is a message in hypertext mark-up language ( HTML ). This is a byte stream of ASCII characters.
- [ ] The message is terminated by the closing of the connection by the server.
- [ ] Well-behaved clients will read the entire document as fast as possible. The client shall not wait for user action (output paging for example) before reading the whole of the document. The server may impose a timeout of the order of 15 seconds on inactivity.
- [x] Error responses are supplied in human-readable text in HTML syntax. There is no way to distinguish an error response from a satisfactory response except for the content of the text.
- [x] The TCP-IP connection is broken by the server when the whole document has been transferred.
- [x] The client may abort the transfer by breaking the connection before this, in which case the server shall not record any error condition.
- [x] Requests are idempotent . The server need not store any information about the request after disconnection.

## References

- The Original HTTP as defined in 1991 - https://www.w3.org/Protocols/HTTP/AsImplemented
- https://http.dev/0.9
