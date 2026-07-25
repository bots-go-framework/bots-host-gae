package gae

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestNewTelegramRedirectTransport_nilBaseURL(t *testing.T) {
	if _, err := NewTelegramRedirectTransport(nil, nil); err == nil {
		t.Error("expected an error for a nil baseURL")
	}
}

func TestNewTelegramRedirectTransport_incompleteBaseURL(t *testing.T) {
	incomplete, err := url.Parse("/no-host-or-scheme")
	if err != nil {
		t.Fatalf("url.Parse: %v", err)
	}
	if _, err = NewTelegramRedirectTransport(incomplete, nil); err == nil {
		t.Error("expected an error for a baseURL without a scheme and host")
	}
}

func TestNewTelegramRedirectTransport_defaultsBaseTransport(t *testing.T) {
	baseURL, err := url.Parse("http://127.0.0.1:1")
	if err != nil {
		t.Fatalf("url.Parse: %v", err)
	}
	transport, err := NewTelegramRedirectTransport(baseURL, nil)
	if err != nil {
		t.Fatalf("NewTelegramRedirectTransport: %v", err)
	}
	if transport.base == nil {
		t.Error("expected a default base RoundTripper when base is nil")
	}
}

func TestTelegramRedirectTransport_rejectsNonTelegramHost(t *testing.T) {
	baseURL, err := url.Parse("http://127.0.0.1:1")
	if err != nil {
		t.Fatalf("url.Parse: %v", err)
	}
	transport, err := NewTelegramRedirectTransport(baseURL, http.DefaultTransport)
	if err != nil {
		t.Fatalf("NewTelegramRedirectTransport: %v", err)
	}
	request, err := http.NewRequest(http.MethodGet, "https://evil.example.com/steal", nil)
	if err != nil {
		t.Fatalf("http.NewRequest: %v", err)
	}
	if _, err = transport.RoundTrip(request); err == nil {
		t.Fatal("expected RoundTrip to reject a non-Telegram destination")
	}
}

func TestTelegramRedirectTransport_rewritesTelegramHost(t *testing.T) {
	var gotPath, gotBody string
	var served bool
	fake := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		served = true
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		gotBody = string(body)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer fake.Close()

	baseURL, err := url.Parse(fake.URL)
	if err != nil {
		t.Fatalf("url.Parse: %v", err)
	}
	transport, err := NewTelegramRedirectTransport(baseURL, nil)
	if err != nil {
		t.Fatalf("NewTelegramRedirectTransport: %v", err)
	}

	request, err := http.NewRequest(
		http.MethodPost,
		"https://"+TelegramAPIHost+"/botTOKEN/sendMessage",
		strings.NewReader(`{"chat_id":1}`),
	)
	if err != nil {
		t.Fatalf("http.NewRequest: %v", err)
	}
	response, err := transport.RoundTrip(request)
	if err != nil {
		t.Fatalf("RoundTrip: %v", err)
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", response.StatusCode, http.StatusOK)
	}
	if gotPath != "/botTOKEN/sendMessage" {
		t.Errorf("path = %q, want /botTOKEN/sendMessage", gotPath)
	}
	if gotBody != `{"chat_id":1}` {
		t.Errorf("body = %q, want the original request body", gotBody)
	}
	if !served {
		t.Fatal("request never reached the fake server — it was not redirected")
	}

	// The original request must be left untouched (RoundTrip clones it).
	if request.URL.Host != TelegramAPIHost {
		t.Errorf("original request was mutated: Host = %q, want %q", request.URL.Host, TelegramAPIHost)
	}
}

func TestBotHostWithHTTPClient_nilClientPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected BotHostWithHTTPClient(nil) to panic")
		}
	}()
	BotHostWithHTTPClient(nil)
}

func TestBotHostWithHTTPClient_returnsSuppliedClient(t *testing.T) {
	client := &http.Client{}
	host := BotHostWithHTTPClient(client)
	if got := host.GetHTTPClient(context.Background()); got != client {
		t.Errorf("GetHTTPClient() = %p, want the supplied client %p", got, client)
	}
}

func TestBotHostWithHTTPClient_GetHTTPClient_nilContextPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected GetHTTPClient(nil) to panic")
		}
	}()
	host := BotHostWithHTTPClient(&http.Client{})
	var ctx context.Context
	host.GetHTTPClient(ctx)
}

func TestBotHostWithHTTPClient_Context(t *testing.T) {
	host := BotHostWithHTTPClient(&http.Client{})
	r := &http.Request{}
	if ctx := host.Context(r); ctx == nil {
		t.Error("Context() returns nil")
	}
}

// TestBotHost_unaffectedByBotHostWithHTTPClient guards the non-breaking
// requirement: BotHost() must keep returning http.DefaultClient regardless of
// BotHostWithHTTPClient's existence.
func TestBotHost_unaffectedByBotHostWithHTTPClient(t *testing.T) {
	_ = BotHostWithHTTPClient(&http.Client{})
	host := BotHost()
	if got := host.GetHTTPClient(context.Background()); got != http.DefaultClient {
		t.Errorf("BotHost().GetHTTPClient() = %p, want http.DefaultClient %p", got, http.DefaultClient)
	}
}
