package gae

import (
	"context"
	"net/http"
	"testing"
)

func TestGaeBotHost_GetHTTPClient_nilRequest(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("The code did not panic")
		}
	}()

	var ctx context.Context = nil
	botHost{}.GetHTTPClient(ctx)
}

func TestBotHost(t *testing.T) {
	v := BotHost()
	if v == nil {
		t.Fatalf("BotHost() returns nil")
	}
}

func TestBotHost_Context(t *testing.T) {
	r := &http.Request{}
	ctx := botHost{}.Context(r)
	if ctx == nil {
		t.Error("Context() returns nil")
	}
}

func TestBotHost_GetHTTPClient(t *testing.T) {
	ctx := context.Background()
	httpClient := botHost{}.GetHTTPClient(ctx)
	if httpClient == nil {
		t.Error("GetHTTPClient() returns nil")
	}
}
