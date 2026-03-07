# ADR-0001: Disabling Timeout Variables 

**DATE:** 05/03/2026 (DD/MM/YYYY)
**Decider:** @ilkerciblak
**Related Issue:** #5


## Context

_Server Sent Event_ communication protocols introduces real-time communication with single and persistence connection over HTTP/HTTPs. SSE keeps the HTTP connection for  a long period while the server continuously sends events or periodic `keep-alive` comment messages to the client. 

The HTTP server was originally configured with standard timeout settings including `WriteTimeout` and `ReadTimeout` which are generally recommended to protect servers from sslow or stalled clients. However this timeout settings also terminates the SSE's persistance HTTP connection between server and client.

## Decision

The HTTP server configuration will disable the `WriteTimeout` and `ReadTimeout` for endpoints that deliver Server-Sent Events (SSE). Since SSE relies on long-lived streaming responses, enforcing a write timeout can prematurely terminate valid connections. Other timeout settings such as `IdleTimeout` will remain enabled to maintain general server protections.


## Consequences

**Positive:**:
- SSE connections may remain open for long periods over single and persistent http connections

**Negative:**
- Resource protection over RESTful endpoints still requires `WriteTimeout` declaration. It can be done using `context` on handlers.
