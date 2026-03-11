# ADR-0002: Refactoring Hub Orchestrator Logic

**DATE**: 10/03/2026
**DECIDER**: [@ilkerciblak](https://github.com/ilkerciblak)
**RELATED ISSUE**: [#10](https://github.com/ilkerciblak/socketoid/issues/10)

## Context

### Presence feature and user-metadata sharing
Currently `Hub` instances orchestrating *SSE* connections and message broadcasting by storing connected clients' `identifiers` and their associated `message channels`.  

Although it was sufficient to manage client connections and broadcast shared messages through all connected clients, introducing `Presence` feature requires to broadcast new three event messages that containing extra user metadata also. By means of `presence feature`, users should be able to retrieve a list of `connected users` with basic user metaadata such as `firstname` or `username`.

Since the current `Client` model does not include any form of **user metadata** the application is unable to expose this information through event messages. As a result, the existing design cannot support the requirements introduced by presence feature.

## Decision

`Client` model is refactored to be have an `Username string` parameter in order to store connected user's name:

```go
type Client struct {
	ID       string      `json:"id"`
	Username string      `json:"user_name"`
	Channel  chan string `json:"-"`
}

```

Also `hub` model now holds `connections` as `map[string]Client` type instead of `map[string]chan string`. Previously it was corresponds to `key: client_id, value: client's message channel`.

```go
type hub struct {
	mu          *sync.Mutex
	connections map[string]Client
	register    chan *Client
	disconnect  chan *Client
	broadcast   chan string
}
```
With this change the `hub` can maintain both the communication channel and the associated client metadata.


## Rationale


**1. Minimal architectural overhead** 
Storing user metadata in a external persiistence layer (e.g. a database) would require additional components, including data access logic and mapping between connection identifiers and stored user records.

Despite this approach may be *preferable* in a production-grade system, it introduces un-necessary architectural complexity beyond the current scope of the application.

**2. Simpler implementation** 
Extending `Client` model to include basic user metadata allows the system to support the required presence functionality with minimal changes. This approach provides a straightforward implementation while preserving the current broadcasting workflow.


## Alternatives Considered

**External user persistence, database or other storage units** 

Maintaining user metadata in a dedicated persistence layer would allow richer user profiles and stronger data consistency guarantees. However this approach would introduce additional layers and development overhead that extends current application's scope.

## Consequences

**Positive**:
- Any event can now broadcast connected user's extra metadata e.g. `name`
- Basic approach to understand and implement

**Negative**:
- Not a production ready development
- `name` field is a little redundant in `Client` model

## Relative ADRs
- [ADR #0003: useSSE Hook Refactoring](./ADR-0003-useSSE-Hook-Refactoring.md)

