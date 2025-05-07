module github.com/NateScarlet/gqlgen-batching/example

go 1.23.4

replace (
	github.com/NateScarlet/gqlgen-batching => ../.
	github.com/NateScarlet/gqlgen-batching/example/generated => ./generated
)

require (
	github.com/99designs/gqlgen v0.17.40
	github.com/NateScarlet/gqlgen-batching v0.0.0-00010101000000-000000000000
	github.com/vektah/dataloaden v0.3.0
	github.com/vektah/gqlparser/v2 v2.5.10
	golang.org/x/text v0.14.0
)

require (
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.3 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sosodev/duration v1.1.0 // indirect
	golang.org/x/mod v0.10.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/tools v0.9.3 // indirect
)
