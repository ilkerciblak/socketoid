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

Every websocket connection is initiated with the _handshake_ request which is a pretty standart `HTTP` request directed from the client party that carries `UPGRADE` headers. Persistent and bidirectional websocket connection is established whether server party responds with `101 Switching Protocol` http respond.

After _websocket_ connection is established, or `handshake process is succeed`, regular HTTP status codes cannot be used to indicate request status. Instead [RFC Spec - Section 7.4](https://www.rfc-editor.org/rfc/rfc6455#section-7.4) states different set of codes.

On the other hand, since the [RFC Spec](https://www.rfc-editor.org/rfc/rfc6455) states, connection should be closed due to several cases mainly violates connection policies. However, in order to not to exceed this application scopes our websocket server will return an `error event`. See [implementation detail](#implementing-error-events) for relevant decision details.

---

## Implementation Details

### Introducing Typed Events on Backend

The communication protocols as _HTTP or Sockets_ aims to introduce rules over distributed system communication. Without these protocols, communication will become vauge and complex to understand.

As a matter of fact, using `typed events` introducing another layer of communication formatting which prevents inconsistent and ambiguous message structures for both parties. In other words, `typed events` introducing a **messaging protocol** by means of providing **consistent structure** to understand and to process for each party.

```go
type event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
```

Supposing the message type is not an unknown message, this structure guarentees:

- Frontend and backend communication over same **contract**. For an example, 'board.card.created' will always be carrying `CardCreatedPayload` to decode.
- Using the `type` field, request will routed to correct handler.

Or if the `type` field corresponds to an unknown event, `readPump` can short-cut the communication with sending `close` frame with `1008 Policy Violation`.

**Using json.RawMessage type over `any`**

In `Go`, unmarshalling `any` fields will return a `map[string]interface{}` type of variable. It would require manual type assertions when application needs `typed payloads`. It will prevent generic `json` operations and also it is an `error prone` operation.

On the other hand using `json.RawMessage` will already stores the value in `bytes` after the parsing operation. It introduces straightforward process while `unmarshalling` wiith using the correct payload struct.

```go
var payload CardCreatedPayload
json.Unmarshal(event.Payload, &payload)
```

**Defined Payload Contracts**

```go
const (
	EventCardCreated = "board.card.created"
	EventCardMoved   = "board.card.moved"
	EventCardDeleted = "board.card.deleted"
	EventCardUpdated = "board.card.updated"
)

type CardCreatedPayload struct {
	CardID string `json:"card_id"`
	Title  string `json:"title"`
	Column string `json:"column"`
}

type CardMovedPayload struct {
	CardID string `json:"card_id"`
	Column string `json:"column"`
}

type CardDeletedPayload struct {
	CardID string `json:"card_id"`
}

type CardUpdatePayload struct {
	CardID string `json:"card_id"`
	Title  string `json:"title"`
	Column string `json:"column"`
}

type errorRespond struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
```

### Orchestrating Event Routing

Previously, incoming messages are handled trivially directly in each client's `readPump` by means of writing back the same incoming message data. As soon as the feature domains are added to applications, mentioned `typed event`'s data should be processed accordingly to the message `type`. In order to route the incoming event to relevant **service**, application requires a centralized `routing` orchestrator that we can `register` services by `type` tagging and `route` these typed event requests.

In `router.go`, orchestrator configured as follows:

```go
type HandlerFunc func(client *client, payload json.RawMessage) error

type router struct {
	routes map[string]HandlerFunc
}

```

`router` instances holds a reference to the registered `type`-`HandlerFunc` couples. Using this hashmap reference, incoming messages can be straightforwardly routed to registered `HandlerFunc`'s using `message type`.

**Helper methods, `Register` & `Route`**

```go
func (r *router) Register(eventType string, handler HandlerFunc) error {
	if _, exists := r.routes[eventType]; exists {
		return fmt.Errorf("event handler already registered as: %v", eventType)
	}

	r.routes[eventType] = handler

	return nil

}

func (r *router) Route(client *client, event event) error {
	handler, exists := r.routes[event.Type]
	if !exists {
		return fmt.Errorf("event handler not registered: %v", event.Type)
	}

	if err := handler(client, event.Payload); err != nil {
		return err
	}

	return nil
}
```

### Implementing Router Logic

Since current connection flow is able to orchestrate messages by its type, client's `readPump` method can use this orchestrator to route the incoming message by its type instead of directly processing it. Using this approach provide more atomic, solid code structure.

```go
func (c *client) readPump(h *hub) {

    for {
//.. other logic
		if opcode == opcodeUTF8Text {
			event, err := UnmarshallEvent(payload)
			if err != nil {
				WriteCloseFrame(c.BuffRW)
				c.cleanUp(h)
				return
			}
			if err := h.router.Route(
				c,
				*event,
			); err != nil {
				errEvent := UnkownEventRespond(event.Type)
				data, e := errEvent.Marshal()
				if e != nil {
					WriteCloseFrame(c.BuffRW)
					return
				}
				WriteFrame(c.BuffRW, data)
				continue
			}

			data, _ := event.Marshal()
			WriteFrame(c.BuffRW, data)
		}

	}

}


```

Although [RFC Spec](https://www.rfc-editor.org/rfc/rfc6455#section-7.4) states, connection should be closed after receiving an `unknown event type`, application responds with a structured `error event`.See [implementation decision](#implementing-error-events).

On the other hand, previous `Hub` signature was also refactored as it is holding a reference to the route orchestrator instance.

```go
type hub struct {
	mu          sync.RWMutex
	register    chan *client
	disconnect  chan *client
	broadcast   chan []byte
	connections map[string]*client
	router      *router
}
```

- Since `Hub` is the main connection orchestrator of the communication server, it is harmonious to hold a `router` reference to its own role.
- This way, a client can share the **same `router` instance** with all connected clients.
- Also, client's dependency will remain same since it can retrieve `router` instance through already injected `hub` instance.

### Introducing Typed Events to Client Party

Similar to [backend part](#introducing-typed-events-on-backend), `typed events` introduce consistent and structured **message protocol** between client and the server. As matter of fact, on both parties these typed events and payloads are defined. On the frontend side, corresponding to backend contracts following definitions were maden:
```ts
eexport interface BoardEvent {
  type: string;
  payload: unknown;
}

export type CardCreatedPayload = {
  Title: string;
  column: string;
};

export class EventCardCreated implements BoardEvent {
  type: string;
  payload: CardCreatedPayload;

  constructor(payload: CardCreatedPayload) {
    this.type = "board.card.created";
    this.payload = payload;
  }
}

export type CardDeletedPayload = {
  card_id: string;
};
export class EventCardDeleted implements BoardEvent {
  type: string;
  payload: CardDeletedPayload;

  constructor(payload: CardDeletedPayload) {
    this.type = "board.card.deleted";
    this.payload = payload;
  }
}

export type CardUpdatedPayload = {
  title: string;
  column: string;
  card_id: string;
};
export class EventCardUpdated implements BoardEvent {
  type: string;
  payload: CardUpdatedPayload;

  constructor(payload: CardUpdatedPayload) {
    this.type = "board.card.updated";
    this.payload = payload;
  }
}

export type CardMovedEventPayload = {
  card_id: string;
  column: string;
};

export class EventCardMoved implements BoardEvent {
  type: string;
  payload: unknown;

  constructor(payload: CardMovedEventPayload) {
    this.type = "board.card.moved";
    this.payload = payload;
  }
}

export const errorEventHandler = (event: MessageEvent): void | Error => {
  const parsed = JSON.parse(event.data);

  if (parsed.type == "error") {
    console.log(JSON.parse(event.data));
    return;
  }

  return new Error(`mis-type used in error-handler: ${event.type}`);
}xport interface BoardEvent {
  type: string;
  payload: unknown;
}

export type CardCreatedPayload = {
  Title: string;
  column: string;
};

export class EventCardCreated implements BoardEvent {
  type: string;
  payload: CardCreatedPayload;

  constructor(payload: CardCreatedPayload) {
    this.type = "board.card.created";
    this.payload = payload;
  }
}

export type CardDeletedPayload = {
  card_id: string;
};
export class EventCardDeleted implements BoardEvent {
  type: string;
  payload: CardDeletedPayload;

  constructor(payload: CardDeletedPayload) {
    this.type = "board.card.deleted";
    this.payload = payload;
  }
}

export type CardUpdatedPayload = {
  title: string;
  column: string;
  card_id: string;
};
export class EventCardUpdated implements BoardEvent {
  type: string;
  payload: CardUpdatedPayload;

  constructor(payload: CardUpdatedPayload) {
    this.type = "board.card.updated";
    this.payload = payload;
  }
}

export type CardMovedEventPayload = {
  card_id: string;
  column: string;
};

export class EventCardMoved implements BoardEvent {
  type: string;
  payload: unknown;

  constructor(payload: CardMovedEventPayload) {
    this.type = "board.card.moved";
    this.payload = payload;
  }
}

export const errorEventHandler = (event: MessageEvent): void | Error => {
  const parsed = JSON.parse(event.data);

  if (parsed.type == "error") {
    console.log(JSON.parse(event.data));
    return;
  }

  return new Error(`mis-type used in error-handler: ${event.type}`);
}

```

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

- https://www.typescriptlang.org/docs/handbook/2/objects.html
- https://hackernoon.com/streaming-in-nextjs-15-websockets-vs-server-sent-events
