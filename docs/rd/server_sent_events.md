> [!Note]
> Developing an application that uses _Server Sent Events_ is almost identically to _websockets_ in part of handling incoming events on the client side. 

## What are Server Sent Events

SSE is a web technology that provides uni-directional real time communication over server to client applications. Using SSE, servers can push data to clients over a _single HTTP connection_. Generally used for the scenarios when real-time updates from the server are essential, but client-to-serveer communication is not necessary

- Key Characteristics
    - Single, persistent HTTP communication
    - Supports uni-directional real-time data flow, from server to client

- Downsides
    - Only supports the MIME type `text/event-stream`

> [!Warning]
> When not used over HTTP/2, SSE suffers from a limitation to the maximum number of open connections, which can be especially painful when opening multiple tabs, as the limit is per browser and is set to a very low number (6). The issue has been marked as "Won't fix" in Chrome and Firefox. This limit is per browser + domain, which means that you can open 6 SSE connections across all of the tabs to www.example1.com and another 6 SSE connections to www.example2.com (per Stack Overflow). When using HTTP/2, the maximum number of simultaneous HTTP streams is negotiated between the server and the client (defaults to 100).

## Event Stream Format

### Content Type and Data Format
The event stream is a simple stream of _text data_ which must be encoded using `UTF-8`. Messages in the event stream are separated by a pair of newline characters `\n`. A colon `:` as the first character of a line s in essence a comment, and is ignored.

- Response header should be set to `text/event-stream`
```go
w.Header().Set("Content-Type", "text/event-stream")
```


### Keeping the Connection Alive
The comment line can be used to _prevent connections from timing out_. A server can send a comment periodically to **keep the connection alive**.
```go
rc := http.NewResponseController(w)
fmt.Fprintf(w, "data: The time is %s\n\n", time.Now().Format(time.UnixDate))
// To make sure that the data is sent immediately
rc.Flush()
```

_CORS and Connection Headers_
```go
w.Header().Set("Cache-Control", "no-cache")
w.Header().Set("Connection", "keep-alive")
w.Header().Set("Access-Control-Allow-Origin", "*")
```

### Fields
Each message received has some combination of the following fields, one per line. In our server we send only the data field which is enough, as other fields are optional. More details here.

**event** – a string identifying the type of event described.

**data** – the data field for the message.

**id** – the event ID to set the EventSource object's last event ID value.

**retry** – the reconnection time.

#### Examples

##### Data-only Messages
```bash
: this is a test stream
```

Presented example introduces a comment since it starts with a colon character. As mentioned previously, this can be useful as a `keep-alive` mechanism if messages might not be sent regularly.

```bash
data: some text
```
Message contains a data field with the value `some text`.

```bash
data: message with
data: two lines
```
Presents a message that contains a data field with the value `"another message\nwith two lines"`. 

##### Named Events
Each stream can have an _event name_ specified by the `event` field, and a `data` field whose value is an appropriate `JSON` string with the data needed for the client to act on the event. PS: the `data` field can be both `string` or `json`

```bash
event: userconnect
data: {"username": "bobby", "time": "02:33:48"}

event: usermessage
data: {"username": "bobby", "time": "02:34:11", "text": "Hi everyone."}

event: userdisconnect
data: {"username": "bobby", "time": "02:34:23"}

event: usermessage
data: {"username": "sean", "time": "02:34:36", "text": "Bye, bobby."}
```

##### Mixing and Matching
You don't have to use just unnamed messages or typed events; you can mix them together in a single event stream.

```bash
event: userconnect
data: {"username": "bobby", "time": "02:33:48"}

data: Here's a system message of some kind that will get used
data: to accomplish some task.

event: usermessage
data: {"username": "bobby", "time": "02:34:11", "text": "Hi everyone."}
```


## Reference
- [MDN's Official Documentation](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events#event_stream_format) for _Using server-sent events_
- _Alex Pliutau_'s [How to Implement Server-Sent Events in Go](https://www.freecodecamp.org/news/how-to-implement-server-sent-events-in-go/) guide on _freecodecamp.org_
