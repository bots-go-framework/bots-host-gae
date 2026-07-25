package gae

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// TelegramAPIHost is the hostname of the real Telegram Bot API
// (github.com/bots-go-framework/bots-api-telegram's tgbotapi.APIEndpoint is
// hardcoded to this host — it has no configurable base URL). It is exported
// so callers can recognize which requests TelegramRedirectTransport rewrites.
const TelegramAPIHost = "api.telegram.org"

// TelegramRedirectTransport is an http.RoundTripper that rewrites requests
// bound for TelegramAPIHost to a caller-supplied base URL and rejects every
// other destination. It exists to let a non-production host (e.g. one wired
// to a Chatwright Telegram Platform Emulator) redirect a bot's outbound Bot
// API calls without turning the process's HTTP client into a general-purpose
// proxy: refusing any request whose host is not TelegramAPIHost is a closed
// network boundary, not an oversight.
//
// Construct one with NewTelegramRedirectTransport and install it on an
// *http.Client passed to BotHostWithHTTPClient.
type TelegramRedirectTransport struct {
	baseURL *url.URL
	base    http.RoundTripper
}

// NewTelegramRedirectTransport returns a TelegramRedirectTransport that
// rewrites the scheme and host of any request whose destination host is
// TelegramAPIHost to baseURL, preserving the path, query and body. Requests
// to any other host are rejected with an error rather than forwarded.
//
// baseURL must be non-nil and specify a scheme and host (e.g.
// "http://127.0.0.1:4000"). base is the underlying RoundTripper used to
// perform the rewritten request; if nil, http.DefaultTransport is used.
func NewTelegramRedirectTransport(baseURL *url.URL, base http.RoundTripper) (*TelegramRedirectTransport, error) {
	if baseURL == nil {
		return nil, errors.New("bots-host-gae: baseURL is nil")
	}
	if baseURL.Scheme == "" || baseURL.Host == "" {
		return nil, fmt.Errorf("bots-host-gae: baseURL %q must have a scheme and a host", baseURL.String())
	}
	if base == nil {
		base = http.DefaultTransport
	}
	return &TelegramRedirectTransport{baseURL: baseURL, base: base}, nil
}

// RoundTrip implements http.RoundTripper. It rejects any request whose host
// is not TelegramAPIHost and otherwise forwards the request, rewritten to
// target the configured baseURL, to the underlying RoundTripper.
func (t *TelegramRedirectTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if request.URL.Host != TelegramAPIHost {
		return nil, fmt.Errorf("bots-host-gae: TelegramRedirectTransport rejected unexpected HTTP destination: %s", request.URL.Host)
	}
	clonedRequest := request.Clone(request.Context())
	clonedURL := *request.URL
	clonedURL.Scheme = t.baseURL.Scheme
	clonedURL.Host = t.baseURL.Host
	clonedRequest.URL = &clonedURL
	return t.base.RoundTrip(clonedRequest)
}
