package gae

import (
	"context"
	"google.golang.org/appengine/v2"
	"net/http"
)

// botHostInterface is the shape returned by BotHost and BotHostWithHTTPClient.
// It intentionally mirrors github.com/bots-go-framework/bots-fw's botsfw.BotHost
// interface without importing that module, so this package keeps its single
// real dependency (google.golang.org/appengine/v2).
type botHostInterface interface {

	// Context returns a context.Context for a request.
	// We need this as some platforms (as Google App Engine Standard)
	// require usage of a context with a specific wrapper
	Context(r *http.Request) context.Context

	// GetHTTPClient returns HTTP client for current host
	// We need this as some platforms (as Google App Engine Standard) require setting http client in a specific way.
	GetHTTPClient(c context.Context) *http.Client
}

// botHost represent information on current hosting platform
type botHost struct {
}

//var _ botsfw.BotHost = (*botHost)(nil)

// BotHost returns hosting platform settings & information. Its GetHTTPClient
// always returns http.DefaultClient; use BotHostWithHTTPClient when a caller
// needs to control outbound HTTP (e.g. to redirect Bot API calls to a
// Chatwright Telegram Platform Emulator in a non-production environment).
func BotHost() botHostInterface {
	return botHost{}
}

// Context creates context for http.Request
func (h botHost) Context(r *http.Request) context.Context {
	return appengine.NewContext(r)
}

// GetHTTPClient creates an HTTP client using AppEngine's URL fetch
func (h botHost) GetHTTPClient(c context.Context) *http.Client {
	if c == nil {
		panic("c == nil")
	}
	return http.DefaultClient
	//return urlfetch.Client(c)
	//return &http.Client{
	//	Transport: &urlfetch.Transport{
	//		Context: c,
	//	},
	//}
}

// botHostWithClient is a botHostInterface whose GetHTTPClient returns a
// caller-supplied *http.Client instead of http.DefaultClient. Context()
// behaves identically to botHost's.
type botHostWithClient struct {
	client *http.Client
}

// BotHostWithHTTPClient returns hosting platform settings & information whose
// GetHTTPClient returns client instead of http.DefaultClient. This is the
// seam a caller uses to redirect a bot's outbound HTTP calls — for example,
// installing a TelegramRedirectTransport to point Bot API traffic at a
// Chatwright emulator instead of https://api.telegram.org — without changing
// anything else about how the host behaves. It does not alter BotHost(),
// which keeps returning http.DefaultClient exactly as before.
//
// client must not be nil; callers that do not need a custom client should
// call BotHost() instead.
func BotHostWithHTTPClient(client *http.Client) botHostInterface {
	if client == nil {
		panic("bots-host-gae: client == nil")
	}
	return botHostWithClient{client: client}
}

// Context creates context for http.Request, identically to botHost.Context.
func (h botHostWithClient) Context(r *http.Request) context.Context {
	return appengine.NewContext(r)
}

// GetHTTPClient returns the *http.Client supplied to BotHostWithHTTPClient.
func (h botHostWithClient) GetHTTPClient(c context.Context) *http.Client {
	if c == nil {
		panic("c == nil")
	}
	return h.client
}

//var DbProvider = func(c context.Context) (db dal.DB, err error) {
//	panic("gae.DbProvider is not set")
//	//return dalgo2datastore.NewDatabase(c, "")
//}

// DB returns database instance
//func (h botHost) DB(c context.Context) (db dal.DB, err error) {
//	if DbProvider == nil {
//		return nil, errors.New("variable DbProvider is not set in github.com/bots-go-framework/bots-host-gae")
//	}
//	return DbProvider(c)
//}

// GetBotCoreStores returns bot DAL
//func (h botHost) GetBotCoreStores(platform string, appContext botsfw.BotAppContext, r *http.Request) (stores botsfwdal.DataAccess) {
//
//	dbProvider := func(c context.Context) (db dal.Database, err error) {
//		return dalgo2datastore.NewDatabase(c, "")
//	}
//
//	//appUserStore := NewGaeAppUserStore(appContext.AppUserEntityKind(), appContext.AppUserEntityType(), appContext.NewBotAppUserEntity)
//	stores.BotAppUserStore = dalgo4botsfw.NewAppUserStore(appContext.AppUserEntityKind(), dbProvider)
//
//	switch platform { // TODO: Should not be hardcoded
//	case "telegram": // pass
//		stores.BotUserStore = dalgo4botsfw.NewBotUserStore("TgUser", dbProvider, func() botsfw.BotUser { return nil }, func(c context.Context, botID string, apiUser botsfw.WebhookActor) (botsfw.BotUser, error) {
//			panic("not implemented")
//		})
//		//if tgChatStore := appContext.GetBotChatEntityFactory(platform); tgChatStore != nil {
//		//	stores.BotChatStore = NewGaeTelegramChatStore(tgChatStore)
//		//} else {
//		//	stores.BotChatStore = NewGaeTelegramChatStore(func() botsfw.BotChat { return telegram.NewTelegramChatEntity() })
//		//}
//		//stores.BotUserStore = newGaeTelegramUserStore(appUserStore)
//	case "fbm": // pass
//		stores.BotUserStore = dalgo4botsfw.NewBotUserStore("FbUser", dbProvider, func() botsfw.BotUser { return nil }, func(c context.Context, botID string, apiUser botsfw.WebhookActor) (botsfw.BotUser, error) {
//			panic("not implemented")
//		})
//		//stores.BotChatStore = NewGaeFbmChatStore()
//		//stores.BotUserStore = newGaeFacebookUserStore(appUserStore)
//	case "viber": // pass
//		stores.BotUserStore = dalgo4botsfw.NewBotUserStore("ViberUser", dbProvider, func() botsfw.BotUser { return nil }, func(c context.Context, botID string, apiUser botsfw.WebhookActor) (botsfw.BotUser, error) {
//			panic("not implemented")
//		})
//		//userChatStore := newGaeViberUserChatStore(appUserStore)
//		//stores.BotChatStore = userChatStore
//		//stores.BotUserStore = userChatStore
//	default:
//		panic("Unknown platform: " + platform)
//	}
//	return
//}
