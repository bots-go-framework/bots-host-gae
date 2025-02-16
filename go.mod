module github.com/bots-go-framework/bots-host-gae

go 1.22.3

//replace github.com/bots-go-framework/bots-fw => ../../bots-go-framework/bots-fw
//
//replace github.com/bots-go-framework/bots-fw-store => ../../bots-go-framework/bots-fw-store

require google.golang.org/appengine/v2 v2.0.6

require (
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)
