
# [Issue #16](https://github.com/ilkerciblak/socketoid/issues/16): WebSocket Event Protocol — Board Event Types & Routing


## Overview

In order to fully discover `real-time communication` over `websockets`, our application needs a feature implementation that will be proccessed over related technologies. As an next step, [Issue #16](https://github.com/ilkerciblak/socketoid/issues/16) subjects _board events_ feature which all CRUD operations will exchanged over _WebSocket_ to conform typed event structures. In order to accomplish this idea process will cover:

- Defining typed events named as `board.card.created`, `board.card.deleted` etc. and corresponding payload structs on backend party.
- Implementing an event router that routes incoming events to related handler function based on `event.Type` field.
- Refactoring `Client` model's `readPump` logic in order to route incoming requests rather then processing requests directly.
- Defining corresponding typed events and payload interfaces in frontend party using `Typescript`
- Exposing `send` wrapper for typed events on frontend party.

---

## Research and Planning

### Error Handling

#### Error Handling on WebSocket Servers

Every websocket connection is initiated with the _handshake_ request which is a pretty standart `HTTP` request directed from the client party that carries `UPGRADE` headers. Persistent and bidirectional websocket connection is established whether server party responds with `101 Switching Protocol` http respond. 

After _websocket_ connection is established, or `handshake process is succeed`, regular HTTP status codes cannot be used to indicate request status. Instead [RFC Spec - Section 7.4](https://www.rfc-editor.org/rfc/rfc6455#section-7.4) states different set of codes.




---


## Implementation Details


---

## Implementation Decisions

### Implementing Error Events

As the [RFC Spec](https://www.rfc-editor.org/rfc/rfc6455#section-7.4) states, unknown event requests violates the connection policy.Thus connection should be closed after receiving an `unknown event type` with `1008 Policy Violation` status code.

However, our project roadmap will cover `feature implementation` in the following steps. Due to that, until this part of the socket implementation unkown events only respond with typed `Error event` similar to the following:

```json
{
    "type": "error",
    "payload": {
        "code": 1008,
        "message": "unknown event type: board.card.foo"
    }
}
```


---

## Related ADRs

---

## References
