# ADR-0003: Refactoring `useSSE` Hook Connection

**DATE**: 10/03/2026
**DECIDER**: [@ilkerciblak](https://github.com/ilkerciblak)
**RELATED ISSUE**: [#10](https://github.com/ilkerciblak/socketoid/issues/10)

## Context

Previous client side _sse_ implementation mechanism was not requiring any type of additional `url parameter` in order to use. Thus using _automatic connection_ with `GET /events` endpoint was sufficient to retrieve and use `event messages`.

Through `presence feature`, connected users should be able to retrieve a list of other connected users with basic user metadata. Since previous fetching mechanism does not required users to fill a prefered name, the application was not able to get this information.

As a result, existing connection mechanism does not provides a system to users to post their usernames and send their prefered name to server.

## Decision

`useSSE` hook's signature was refactored as follows:

```typescript
export default function useSSE(
  url: string,
  options?: {
    handlers: onMessage[];
    auto: boolean;
  },
);
```

Following this change, `connect` functions' signature is changed also to have an _optional_ `query` parameter. Later on, using this optional query parameter and basic `URL()` interface api, _sse_ request url is structured.

```typescript
const connect = useCallback(
    (query?: { key: string; val: string }) => {
      if (eventSourceRef.current) eventSourceRef.current.close();
      let base: URL = new URL(url);
      if (query) {
        base.searchParams.set(query.key, query.val);
      }
      const eventSource = new EventSource(base, { withCredentials: false });

//... event handling logics (stayed same)
```

Lastly, previous `useEffect` part was containing no conditional logic to run `connect` method. This was altered to have a conditional as

```typescript
useEffect(() => {
  if (optionsRef.current?.auto) {
    connect();
  }
  return () => {
    eventSource.current?.close();
  };
}, [connect, url]);
```

## Rationale

**1. Encapsulated Reusable Logic**

Using same hook that implements new and prevous application requirements provides extendable re-usable, solid and encapsulated business logic to use.

**2. Basic implementation**

Extending previous custom hook allowed the system to support the required functionality with minimal and straightforward implementations. Also, previous user interfaces (e.g. health page) required no changes in hook implementation.

## Alternatives Considered

**Additional function definition in same hook**

Similar to `connect` function a new function with optional url part parameters that not triggered at `useEffect` part could be declared. This would also provide extension on the same hook with `reusable, testable and encapsulated` interface. However, since both (new connect function and the old one) would require same connection and event handling functionalities, this approach would introduce code smell or development overhead for code restructring.

## Consequences

**Positive**:

- Additional `query` parameter can be instrumented to _SSE_ endpoint
- Connection can be open on prefered point of the application.
- Logic is encapsulated and reusable, also testable.

**Negative**:

## Relative ADRs

-
