# Issue#8 - Server Sent Event Client Integration on Next.js Frontend [See on GitHub repository](https://github.com/ilkerciblak/socketoid/issues/8)


## Overview

## Research

### EventSource Interface 

The `EventSource` interface is web content's interface to _server sent events_. An `EventSource` instance opens a persistance connection to an **HTTP server**. The connection _remains open_ until _closed by calling `EventSource.close()`.

Once the connection is opened, incoming messages from the server are delivered to client n the form of _events_. If there is an `event field` in the incoming message, the triggered event is the same as the event field value. If no event field is present, then a generic `message event` is fired.

Named events mentioned, can be listened using basic `addEventListener('event-name' ()=>{})` usage. [See **Examples** section](#### Examples)

#### Instance Properties

> [!NOTE]
> This interface also inherits properties from its parent `EventTarget`, see [EventTarget on MDN Documentation](https://developer.mozilla.org/en-US/docs/Web/API/EventTarget)

-  `EventSource.readyState` | _readonly_ <br/>
    A number representing the state of the connection. Possible values are `0`, `1` and `2` for `CONNECTING`, `OPEN` and `CLOSED` respectively.

- `EventSource.url` | _readonly_ <br/>
    A string representng the `URL` of the source.

- `EventSource.withCredentials` | _readonly_ <br/>
    A `boolean` value indicating whether the `EventSource` object was instantated with `CORS` credentials set.  

#### Instance Methods

> [!NOTE]
> This interface also inherits properties from its parent `EventTarget`, see [EventTarget on MDN Documentation](https://developer.mozilla.org/en-US/docs/Web/API/EventTarget)  

- `EventSource.close()`: 
    Closes the connection, and sets the `readyState` attribute to `CLOSED` if any.

> [!IMPORTANT]
> If there is an `event` field in a message, client js should use `addEventListener('event-name')` instead of `onmessage`. 

- `EventSource.addEventListener('event-name', f)` | _inherited from `EventTarget`_<br/>
    The addEventListener() method of the EventTarget interface sets up a function that will be called whenever the specified event is delivered to the target. See [official documentation on MDN](https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/addEventListener) if needed.



#### Events
- `error`: Fired when a connection to an event source _failed to open_.
- `message`: Fired when _data is received_ from an event source.
- `open`: Fired when a _connection to an event source is opened_.

#### Examples

#### Un-named Events -> `onmessage`
In this basic example an `EventSource` is created to receive _unnamed events_ from the server:
```js
const evtSource = new EventSource("sse.php");
const eventList = document.querySelector("ul");

evtSource.onmessage = (e) => {
  const newElement = document.createElement("li");

  newElement.textContent = `message: ${e.data}`;
  eventList.appendChild(newElement);
};

```
Each received event causes `onmessage` event handler to be run. It creates a new `<li>` element and writes the message's data into it.


#### NamedEvents -> `addEventListener`

In order to listen named events, it is required to define a listener for each type of event sent:
```javascript
const sse = new EventSource("some-url");

sse.addEventListener("notice", (e) => {
    console.log(e.data)
})
/*
 * This will listen only for events
 * similar to the following:
 *
 * event: notice
 * data: useful data
 * id: some-id
 */

 /*
 * Similarly, this will listen for events
 * with the field `event: update`
 */
sse.addEventListener("update", (e) => {
  console.log(e.data);
});

/*
 * The event "message" is a special case, as it
 * will capture events without an event field
 * as well as events that have the specific type
 * `event: message` It will not trigger on any
 * other event type.
 */
sse.addEventListener("message", (e) => {
  console.log(e.data);
});


```
:
## Implementation Notes

Despite there are ready-made solutions, server sent event's connection lifecycle, event processing and cleanup logic is encapsulated in a custom `useSSE` hook. Although, for now, there are no multiple Sse event usage in this project, encapsulating the logic introduces clean, re-usable and testable code structure.

Our hook signature looks like following:
```typescript
type onMessage = { type: string; handler: (event: MessageEvent) => void };

export default function useSSE(
  url: string,
  options?: {
    handlers: onMessage[];
  },
)
```

Since default `EventSource.onmessage` method does not catch and processes _named events_, user interface can define the `type` of the event and `handler` as what to do with this event message data. For an example, our `SSEStatusIndicator` component tracks `userjoined` event and generates a list of `user-id`s coming from the event data.

```typescript
export default function SSEStatusIndicator() {
  const [userList, setUserList] = useState<string[]>([]);
  const { connectionState, error: errorMsg } = useSSE(
    app_config.PUBLIC_API + "/events",
    {
      handlers: [
        {
          type: "userjoined",
          handler(event) {
            const data = JSON.parse(event.data);
            setUserList((prev) => [data["user-id"], ...prev]);
          },
        },
      ],
    },
  );
```

In addition `useSSE` hook also returns `connectionStatus` and `errorMsg` information to use in the indicator UI. Those are aiming miscellanous.

Mainly to process sse connection lifecycle native `EventSource` interface API is used. Handling `EventSource.onopen` and `EventSource.onerror` events provided easy and clean implementations. On the other hand, `EventSource.onmessage` event is not sufficient for event data processing since its not catching _named events_. Thus hook signature allows user interface to inject some type of event data processing handlers.

```typescript
if (optionsRef?.current?.handlers.length) {
        optionsRef?.current?.handlers.forEach((handler) => {
          eventSource.addEventListener(handler.type, (e) => {
            try {
              handler.handler(e);
            } catch (error) {
              setErrorMsg(`failed to process message: ${error}`);
            }
          });
        });
      } else {
        eventSource.onmessage = (event) => {
          try {
            const data = JSON.parse(event.data);
            setEventQueue((prev) => [...prev, data]);
          } catch (error) {
            setErrorMsg(`${error}`);
          }
        };
      }
```

Lastly since all `React Hooks` are functions, using `useEffect` hook with `connect` and a _clean up function_, is lets automatically binding the connection with the component that this hook is used. `EventSource.close()` instance method is used as clean up function to prevent leaks when the rendering is over.

```typescript
useEffect(() => {
    connect();
    return () => {
        eventSourceRef.current?.close()
    }
}, [connect, url])
```
## Implementation Decisions 
***Rendering Optimizations with `useCallback` and `useRef`***
Using `connect` method as a direct dependency for `useEffect` was resulting infinite re-rendering. 

```
page renders -> hook's useEffect triggered -> connect() called -> new EventSource -> new State -> re-render -> loop starts over
```

Thus `connect` function defined using `useCallback` hook with `url` dependency.

Also `options` parameter was causing same context. Thus in order to solve this issue `useRef` hook was used. 

## Related ADRs

## References

- Official MDN Documentation on Event Source: https://developer.mozilla.org/en-US/docs/Web/API/EventSource
- Official react.dev Documentation on Custom Hooks: https://react.dev/learn/reusing-logic-with-custom-hooks
