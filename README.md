# gqlgen-batching

GraphQL batch support for gqlgen

## What is gqlgen-batching?

**gqlgen-batching** is an extension of [gqlgen](https://github.com/99designs/gqlgen) to support [GraphQL Batching](https://github.com/graphql/graphql-over-http/blob/main/rfcs/Batching.md).

## Features

- Automatically detects batch requests by checking if the payload starts with `[`
- Processes queries in parallel according to configured concurrency settings
- Streams responses as they become available
- Maintains proper status codes (may return 422 for syntax errors while other requests succeed with 200)
- Only handles supported requests without interfering with regular POST transport

## Quick start

1. Add package import

   ```go
   import "github.com/NateScarlet/gqlgen-batching/pkg/batching"
   ```

2. Add batching transport

   ```go
   srv.AddTransport(batching.POST{})
   ```

3. Prepare your server

   ```go
   schema := generated.NewExecutableSchema(graph.NewResolver())
   srv := handler.New(schema)

   srv.AddTransport(transport.Websocket{
   	KeepAlivePingInterval: 10 * time.Second,
   })
   srv.AddTransport(transport.Options{})
   srv.AddTransport(transport.GET{})
   srv.AddTransport(batching.POST{}) // Handles batch requests
   srv.AddTransport(transport.POST{}) // Handles normal requests
   srv.AddTransport(transport.MultipartForm{})

   srv.SetQueryCache(lru.New(1000))

   srv.Use(extension.Introspection{})
   srv.Use(extension.AutomaticPersistedQuery{
   	Cache: lru.New(100),
   })

   http.Handle("/", playground.Handler("Starwars", "/query"))
   http.Handle("/query", srv)

   log.Fatal(http.ListenAndServe(":8080", nil))
   ```

## How to use?

Simply send a JSON array of GraphQL requests:

```bash
curl -X POST http://localhost:8080/query \
-H "Content-Type: application/json" \
-d '[{"query":"{hero(episode: JEDI) { name }}"},{"query":"{hero(episode: EMPIRE) { name }}"}]'
```

Result:

```json
[
  {
    "data": {
      "hero": {
        "name": "R2-D2"
      }
    }
  },
  {
    "data": {
      "hero": {
        "name": "Luke Skywalker"
      }
    }
  }
]
```
