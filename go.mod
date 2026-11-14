module github.com/bots-go-framework/bots-host-gae

go 1.26.0

toolchain go1.27.4

//replace github.com/bots-go-framework/bots-fw => ../../bots-go-framework/bots-fw
//
//replace github.com/bots-go-framework/bots-fw-store => ../../bots-go-framework/bots-fw-store

require google.golang.org/appengine/v2 v2.1.0

require (
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
