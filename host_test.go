package gae

import (
	"context"
	"testing"
)

const noPanic = "The code did not panic"

func TestGaeBotHost_GetHTTPClient_nilRequest(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error(noPanic)
		}
	}()

	var ctx context.Context = nil
	botHost{}.GetHTTPClient(ctx)
}
