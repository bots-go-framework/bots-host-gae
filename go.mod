module github.com/bots-go-framework/bots-host-gae

go 1.21

toolchain go1.22.6

//replace github.com/bots-go-framework/bots-fw => ../../bots-go-framework/bots-fw
//
//replace github.com/bots-go-framework/bots-fw-store => ../../bots-go-framework/bots-fw-store

require (
	github.com/bots-go-framework/bots-fw v0.25.2
	github.com/dal-go/dalgo v0.12.1
	google.golang.org/appengine/v2 v2.0.6
)

require (
	github.com/bots-go-framework/bots-fw-store v0.4.0 // indirect
	github.com/bots-go-framework/bots-go-core v0.0.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/pquerna/ffjson v0.0.0-20190930134022-aa0246cd15f7 // indirect
	github.com/strongo/gamp v0.0.1 // indirect
	github.com/strongo/i18n v0.0.4 // indirect
	github.com/strongo/random v0.0.1 // indirect
	github.com/strongo/validation v0.0.6 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)
