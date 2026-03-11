# Issue#10 - Presence System Instrumentation [See on GitHub repository](https://github.com/ilkerciblak/socketoid/issues/10)

## Overview

User presence systems track connected users and broadcasts _join/disconnect_ events to all connected clients. When a user connects to an application's server, server notifies other clients with `user.joined` like events. When a user disconnects, the same applies.

The `presence feature` requires an unidirectional, server driven data flow to clients. Consequently this feature is the most natural use case for SSE client and server implementation.

## Research and Planning

No research.

### Functional Requirements

- Connected users should be able to other connected users' name/usernames.
- Users should be able to fill a prefered name before connection
- Every time an user connect or disconnected from the server, the connected user list should react the new state.

## Implementation Notes

In order to conduct client notifications through previous message broadcasting mechanism, three new `user.joined`, `user.left` and `presence.init` events were required. These event messages requires to provide connected clients' information with basic _user metadata_.

- `user.joined`: Event message type maintains a single user's metadata information that joins to connection
- `user.left`: Event message type corresponds to client that disconnect from the server connection, maintains a single user's metadata
- `presence.init`: Event message type maintains a _online client_ list with basic _user metadata_ information.

To provide solid structure for _SSE Event_ and _presence event_ signatures, following `PresencePayload` and `SSEEvent` models are defined.

```go
type PresencePayload struct {
	UserId string `json:"user_id"`
	Name   string `json:"name"`
}

type SSEEvent struct {
	Type    string `json:"type"` //name of the event e.g. user.joined
	Payload any    `json:"payload"` // payload of the event message, data part
}
```

`Payload` parameter type can be defined as a `SSEPayloadInterface` later. However it is sufficient use as `any` for now. 

In addition to model definitions, public functions are declared to construct concrete event structuring. This function signatures introducing clear event creation interfaces while only accepts payload parameter:

```go
const (
	USER_JOINED string = "user.joined"
	USER_LEFT   string = "user.left"
)

func UserJoinedEvent(payload PresencePayload) SSEEvent {
	return SSEEvent{
		Type:    USER_JOINED,
		Payload: payload,
	}
}

func UserLeftEvent(payload PresencePayload) SSEEvent {
	return SSEEvent{
		Type:    USER_LEFT,
		Payload: payload,
	}
}
```

As mentioned, one downside of the _SSE endpoints_ is they only supports the MIME type `text/event-stream`. Thus all events should be transformed into some structured `text` type. Thus, to introduce reusable mechanism `SSEEvent.ToTextStream() string` method is declared as follows:

```go
func (e SSEEvent) ToTextStream() string {
	dataByte, err := json.Marshal(e.Payload)
	if err != nil {
		return fmt.Sprintf("event: error\ndata: %s", err.Error())

	}
	data := fmt.Sprintf("data: %s", dataByte)
	event := fmt.Sprintf("event: %s", e.Type)
	return event + "\n" + data + "\n\n"

}
```

Using this method, any `SSEEvent` instance can be mapped into strucuted event-stream text. For an example, `user.joined` event was structured and broadcast in _client registration_ process as follows:

```go
case new_client := <-h.register:
	if err := h.registerClient(*new_client); err != nil {
		fmt.Println(err)
	} //client registration 
    // Creating payload instance
	payload := PresencePayload{
		Name:   new_client.Username,
		UserId: new_client.ID,
	}

	event := UserJoinedEvent(payload)
    // Sending connected user list to just registered client
	new_client.Channel <- h.ConnectionListEvent(new_client.ID)
    // broadcasting `user.joined` event to other clients
	h.broadcastMessage(event.ToTextStream())
```

On the other hand, connected users should retrieve an event message that maintains _online user list_ when they connect to our _sse server_. Since `hub` orchestrates these connections and connected users, processing this event handling in `hub` logic was introducing simpler architecture and implementation.

```go
func (h *hub) ConnectionListEvent(except string) string {
	h.mu.Lock()
	var userList []PresencePayload
	for _, client := range h.connections {
		if !strings.EqualFold(except, client.ID) {
			userList = append(userList, PresencePayload{
				UserId: client.ID,
				Name:   client.Username,
			})
		}
	}
	h.mu.Unlock()

	dataByte, _ := json.Marshal(userList)
	return fmt.Sprintf("event: presence.init\ndata: %s\n\n", string(dataByte))

}
```

Similar to `ToTextStream` method, `Hub.ConnectionListEvent` method structures a acceptable `text/event-stream` type of string to broadcast it to connected client. It produces `presence.init` event with _json string_ type of data that contains list of connected users' metadata.

## Implementation Decisions

**1. Structuring `presence.init` event message in hub orchestrator**

Initial event stream of the presence feature requires to maintain a list of connected user metadata. Since hub orchestrates the server sent event connections and connected users, it introduces less complexity and sufficient functionality to structure this message as a method of hub orchestrator.

**2. Refactoring SSE Handler and Hub.Run method**

Handler's responsibility is solely to manage the HTTP connection and introduce the client to the Hub — not to make business logic decisions. When the user.joined broadcast and presence.init delivery were handled inside the handler, it was making decisions about "what should happen after registration", which coupled the handler to the Hub's internal behavior. By moving this logic into the Run() loop, all these decisions are centralized within the Hub itself — registration becomes atomic, the handler simply writes h.register <- &client and its job is done, while the Hub coordinates everything else.

## Related ADRs

-See [related ADR documentation](../adr/ADR-0002-Refactoring-Hub-Orchestrator-Logic.md)

- On SSE client implementation, due to connection logic changes `useSSE` hook is refactored by means of adding boolean parameters. See [related ADR documentation](../adr/ADR-0003-useSSE-Hook-Refactoring.md)
   

## References

- See [CodeSignal Blog Post on JSON Encoding/Decoding](https://codesignal.com/learn/courses/handling-json-in-go-1/lessons/encoding-structs-into-json-in-go-1)
