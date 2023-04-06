package gae

import (
	"context"
	"github.com/bots-go-framework/bots-fw-dalgo/dalgo4botsfw"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo2datastore"
	"google.golang.org/appengine/v2"
	"google.golang.org/appengine/v2/urlfetch"
	"net/http"
)

// botHost represent information on current hosting platform
type botHost struct {
}

var _ botsfw.BotHost = (*botHost)(nil)

// BotHost returns hosting platform settings & information
func BotHost() botsfw.BotHost {
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
	return &http.Client{
		Transport: &urlfetch.Transport{
			Context: c,
		},
	}
}

// DB returns database instance
func (h botHost) DB() dal.Database {
	panic("not implemented")
	//return gaedb.NewDatabase()
}

// GetBotCoreStores returns bot DAL
func (h botHost) GetBotCoreStores(platform string, appContext botsfw.BotAppContext, r *http.Request) (stores botsfw.BotCoreStores) {

	dbProvider := func(c context.Context) (db dal.Database, err error) {
		return dalgo2datastore.NewDatabase(c, "")
	}

	//appUserStore := NewGaeAppUserStore(appContext.AppUserEntityKind(), appContext.AppUserEntityType(), appContext.NewBotAppUserEntity)
	stores.BotAppUserStore = dalgo4botsfw.NewAppUserStore(appContext.AppUserEntityKind(), dbProvider)

	switch platform { // TODO: Should not be hardcoded
	case "telegram": // pass
		stores.BotUserStore = dalgo4botsfw.NewBotUserStore("TgUser", dbProvider, func() botsfw.BotUser { return nil }, func(c context.Context, botID string, apiUser botsfw.WebhookActor) (botsfw.BotUser, error) {
			panic("not implemented")
		})
		//if tgChatStore := appContext.GetBotChatEntityFactory(platform); tgChatStore != nil {
		//	stores.BotChatStore = NewGaeTelegramChatStore(tgChatStore)
		//} else {
		//	stores.BotChatStore = NewGaeTelegramChatStore(func() botsfw.BotChat { return telegram.NewTelegramChatEntity() })
		//}
		//stores.BotUserStore = newGaeTelegramUserStore(appUserStore)
	case "fbm": // pass
		stores.BotUserStore = dalgo4botsfw.NewBotUserStore("FbUser", dbProvider, func() botsfw.BotUser { return nil }, func(c context.Context, botID string, apiUser botsfw.WebhookActor) (botsfw.BotUser, error) {
			panic("not implemented")
		})
		//stores.BotChatStore = NewGaeFbmChatStore()
		//stores.BotUserStore = newGaeFacebookUserStore(appUserStore)
	case "viber": // pass
		stores.BotUserStore = dalgo4botsfw.NewBotUserStore("ViberUser", dbProvider, func() botsfw.BotUser { return nil }, func(c context.Context, botID string, apiUser botsfw.WebhookActor) (botsfw.BotUser, error) {
			panic("not implemented")
		})
		//userChatStore := newGaeViberUserChatStore(appUserStore)
		//stores.BotChatStore = userChatStore
		//stores.BotUserStore = userChatStore
	default:
		panic("Unknown platform: " + platform)
	}
	return
}
