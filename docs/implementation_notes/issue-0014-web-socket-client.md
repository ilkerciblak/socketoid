# Issue#14 - Implementing WebSocket Client [See on GitHub repository](https://github.com/ilkerciblak/socketoid/issues/14)

## Overview

Real time communication protocols referes to continous exchange of data with minimal latency. Unlike traditional request-response model where _data exchange happens on demand_, live updates happens w/o requirement of refreshing UI.

The WebSocket Protocol is a real-time communication protocol that provides bi-directional data exchange, mostly this behavior called full-dublex communication, over persistent single TCP connection between client and the server.

To communicate _websocket_ servers, client parties requires to be instrumented related implementations. Although there are third party solutions provides solid implementations for _websocket client_ instrumentation, browsers native `WebSocket Interface` introduces robust API for creating and managing connection to the server, as well as for sending and receiving data on the connection.

---

## Research and Planning

### Writing WebSocket Client Applications

#### Creating a `WebSocket` object

The `WebSocket` constructor takes one mandatory argument, **the URL** of the websocket server to connect to.

```js
const uri = "ws://localhost:80/ws;
const websocket = new WebSocket(uri);
```

> [!IMPORTANT]
> Similar to HTTP and HTTPs, WebSockets have a unique set of prefixes: **ws** and **wss** for the connections without and with`TLS` respectively. This issue will cover `ws` connections since its connecting `localhost` only. Production level applications should be served using `wss` as the protocol.

The constructor will throw a `SecurityError` if the destination does **not** allow access. This also may happend due to attemps on **insecure connections** due most user agents now require a secure link for all _websocket_ connections unless either party on the _same device_ or _on the same network_.

The `WebSocket` constructor takes another **optional** argument [`protocols`](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket/WebSocket#protocols), a single string or an array of strings, to implement multiple **sub-protocols.**

#### Listening `WebSocket` Events

`WebSocket` Interface API has following events to be able to listen, `close`, `error`, `message` and `open`. These _events_ can be listened using native javascript's `addEventListener()` or by assigning an event listener to the `oneventname` property of this interface.

For an example, once the connection is established, the `open` event is fired. Following code example, sending one _ping_ message to the server every second after connection is opened.

```js
websocket.addEventListener("open", () => {
  log("CONNECTED");
  pingInterval = setInterval(() => {
    log(`SENT: ping: ${counter}`);
    websocket.send("ping");
  }, 1000);
});
```

- `close`

  Fired when a connection with a WebSocket is closed. Also available via the `onclose` property

- `error`

  Fired when a connection with a WebSocket has been closed because of an error, such as when some data couldn't be sent. Also available via the `onerror` property.

- `message`

  Fired when data is received through a WebSocket. Also available via the `onmessage` property.

- `open`

  Fired when a connection with a WebSocket is opened. Also available via the `onopen` property.

#### Sending and Receiving Messages

##### Sending Messages

An `WebSocket` instance can only _send message(s)_ once the connection is established and is alive.

The `websocket.send(data)` method **enqueues** the specified data to be transmitted to server. It increasing the value of [`bufferedAmount`](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket/bufferedAmount) by the number of bytes needed to contain the data.

Whether the data can **not** be sent, buffer might be full or any other error may occur, **the socket closed automatically.** The browser will throw an exception if `send()` is called during `CONNECTING` state. On the other hand, the browser will **silently discards** the data on `CLOSING` or `CLOSED` states of connection.

The `send()` method is **asynchronous**. It does **not** wait for the data to be transmitted before returning to the caller. It just adds the data to its internal buffer and begins the process of transmission. The [`WebSocket.bufferedAmount`](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket/bufferedAmount) property represents **the number of bytes that have not yeet been transmitted**.

> [!IMPORTANT]
> If protocol uses `UTF-8` to encode text, so `bufferedAmount` is calculated based on the `UTF-8` encoding of any buffered text data.

Nevertheless UTF-8 text type of **data** is mostly used, `data` may sent over as `Blob`, `ArrayBuffer`, `TypedArray` or `DataView`. See [MDN Documentation over WebSocket.Send](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket/send#parameters) for more details.

##### Sending JSON Data

A **common approach** to use `JSON` to send **serialized** objects as **text**. For example, instead of just sending the text message "ping", our client could send a serialized object including the number of messages exchanged so far:

```js
const message = {
  iteration: counter,
  content: "ping",
};
websocket.send(JSON.stringify(message));
```

##### Receiving Messages

In order to handle message receiving part, application can listen for the `message` event.

The server can also send binary data, which is exposed to clients as a `Blob` or an `ArrayBuffer` based on the value of the `WebSocket.binaryType` property.

```js
websocket.addEventListener("message", (e) => {
  const message = JSON.parse(e.data);
  log(`RECEIVED: ${message.iteration}: ${message.content}`);
  counter++;
});
```

##### Binary Type Property

The `WebSocket.binaryType` instance property controls the type of binary data being **received** over the websocket connection. It is a string type of property that can be set `blob` which is the default or `arraybuffer`.

```js
// Create WebSocket connection.
const socket = new WebSocket("ws://localhost:8080");

// Change binary type from "blob" to "arraybuffer"
socket.binaryType = "arraybuffer";

// Listen for messages
socket.addEventListener("message", (event) => {
  if (event.data instanceof ArrayBuffer) {
    // binary frame
    const view = new DataView(event.data);
    console.log(view.getInt32(0));
  } else {
    // text frame
    console.log(event.data);
  }
});
```

---

## Implementation Details

[Issue 14](https://github.com/ilkerciblak/socketoid/issues/14) subjects to implement _websocket_ client to the frontend application. Scope and decisions requireds to develop a custom hook that able to connect given _websocket endpoint_, and can handle connection lifecyle and real-time data transmission events. In addition to that, no other third party integrations are restricted despite using browser's native `WebSocket Interface`. A basic _connection status representer_ UI component is also required to test out the connection more straightforward.

### Creating Custom WebSocket Connection Hook

In order to provide re-usable functionalities and easy implementations with a custom hook, `useWebSocket` hook handles socket events implicitly `open, error, close, message` and exposes `send and close` functionalities. However, hook signature also accepts optional event handlers to be run on relevant connection or transmission event.

```ts
export default function useWebSocket(g: {
  url: string;
  messageListeners: onMessage[];
  onOpen?: onEvent;
  onDisconnect?: onEvent;
  onError?: onEvent;
});

// where

type onMessage = {
  type: string;
  handler: (event: MessageEvent) => void;
};

type onEvent = (event: Event) => void;
```

**Initializing State Variables**

```ts
export default function useWebSocket(g: {
  url: string;
  messageListeners: onMessage[];
  onOpen?: onEvent;
  onDisconnect?: onEvent;
  onError?: onEvent;
}) {
  const [connectionState, setConnectionState] = useState<WsConnectionStatus>(
    WsConnectionStatus.IDLE,
  );

  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const webSocketRef = useRef<WebSocket | null>(null);
  const listeners = useRef<onMessage[]>(g.messageListeners);
```

As a first step of the implementation, internal state parameters and references are defined to manage the lifecycle and behavior of the connection:

- `connectionState`: A state variable that represents the current connection status. It is initialized with `WsConnectionStatus.IDLE` which indicates no connection attempts or results yet.
- `errorMsg`: Another straightforward state variable that represents the possible error encountered during connection lifecyle. It is initialized as null and updated when an error occurs, allowing the consumer of the hook to react accordingly (e.g., displaying error feedback in the UI).
- `webSocketRef`: A mutable reference to an active `WebSocket` instance. This reference aiming to prevent additional render triggerations and ensuring using same `WebSocket` instance through the same component's lifecyle.
- `listeners`: A mutable reference containing the array of message listener callbacks (onMessage[]) provided via the hook's input. By storing listeners in a useRef, the hook avoids unnecessary re-subscriptions or re-initializations when the component re-renders, while still allowing access to the latest listener set.

**Establishing Connection**

Websocket connections always initiated from the client party via initial _handshake request_. As mentioned, `useWebSocket` hook implicitly handles connection lifecycles while accepting listener methods in its signature.

On the other hand, creating a `WebSocket` instance starts the process of establishing a connection to the server. Thus in the private `connect` function:

```ts
export default function useWebSocket(g: {
  url: string;
  messageListeners: onMessage[];
  onOpen?: onEvent;
  onDisconnect?: onEvent;
  onError?: onEvent;
}) {
  //.. state vaariable part

  const connect = useCallback(() => {
      // clears the previous reference whether there is
    if (webSocketRef) webSocketRef.current = null;

    try {
      const ws = new WebSocket(g.url);
      setConnectionState(WsConnectionStatus.CONNECTING);
      webSocketRef.current = ws;
      //..code
```

Connection initiating part was wrapped with a `trycatch block` to handle constructer errors. After connection is initiated, connection status is set to `CONNECTING` until the `onopen` event is triggered.

```ts
try {
      const ws = new WebSocket(g.url);
      setConnectionState(WsConnectionStatus.CONNECTING);
      webSocketRef.current = ws;

      //handling onopen event
      webSocketRef.current.onopen = (event) => {
        setConnectionState(WsConnectionStatus.OPEN);
        if (g.onOpen) g.onOpen(event);
      };

//code..
```

`onopen` handler alters the connection status, and optionally calls the defined `onOpen` callback parameter.

Similarly, `connect` methods declares other event handlers implicitly,

```ts
webSocketRef.current.onmessage = (event) => {
  try {
    listeners.current.forEach((listener) => listener.handler(event));
  } catch (e) {
    setErrorMsg(`${e}`);
  }
};

webSocketRef.current.onclose = (event) => {
  setConnectionState(WsConnectionStatus.CLOSED);
  if (g.onDisconnect) g.onDisconnect(event);
};

webSocketRef.current.onerror = (event) => {
  setConnectionState(WsConnectionStatus.ERROR);
  if (g.onError) g.onError(event);
};
```

**close and send handlers**

Hook also defines and exposes data sending and connection closing functionalities.

Since application scope is not requiring _media transfer, high frequency data traffic or custom binary protocol implementations_, using `JSON Stringified` text data for transmisson is sufficient. It will provide straightforward debugging and neglected complexity for data serialization and deserialization.

On the other hand, `WebSocket` interface api **implicitly masking** the payload. Thus any masking logic was used in the `sendMessage` function.

```ts
const sendMessage = useCallback(
  (data: { type: string; payload: unknown }) => {
    try {
      if (
        webSocketRef.current &&
        webSocketRef.current.readyState != webSocketRef.current.OPEN
      ) {
        throw new Error(
          `connection is not open: ${webSocketRef.current.readyState}`,
        );
      }
      const serialized = JSON.stringify(data);
      webSocketRef.current?.send(serialized);
    } catch (error) {
      setErrorMsg("serialization error");
      console.log(`serialization error ${error}`);
    }
  },
  [webSocketRef],
);
```

Exposed `closeConnection` is straightforward.

```ts
const closeConnection = useCallback(() => {
  webSocketRef.current?.close();
}, [webSocketRef]);
```

**Control Frames** 

Modern browsers preventing creating and sending `ping-pong` frames while automatically providing a `keep-alive` functionality themselves.  

**Configuring hook lifecyle** 

Using `React`'s `useEffect` hook, our custom hook's `setup` and `clean-up` logic is configured.

```ts
useEffect(() => {
    connect();

    return () => {
      if (webSocketRef?.current?.readyState != WebSocket.CLOSED) {
        webSocketRef?.current?.close();
      }
    };
  }, [g.url]);
```
Given configuration ensures, each time component commits React will automatically initites ws connection. In contrast, when the component is removed from the DOM React will automatically close the connection. 

---

## Related ADRs

## References

- https://developer.mozilla.org/en-US/docs/Web/API/WebSocket
- https://react.dev/reference/react/useEffect
