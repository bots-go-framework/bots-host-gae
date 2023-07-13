package gae

//
//import (
//	"context"
//	"fmt"
//	"github.com/bots-go-framework/bots-fw/botsfw"
//	"github.com/strongo/log"
//	"github.com/strongo/nds"
//	"google.golang.org/appengine/v2/datastore"
//	"reflect"
//)
//
//// GaeAppUserStore DAL
//type GaeAppUserStore struct {
//	appUserEntityKind string
//	appUserEntityType reflect.Type
//	newUserEntity     func() botsfw.BotAppUser
//	GaeBaseStore
//}
//
//var _ botsfw.BotAppUserStore = (*GaeAppUserStore)(nil)
//
//// NewGaeAppUserStore created new DAL
//func NewGaeAppUserStore(appUserEntityKind string, appUserEntityType reflect.Type, newUserEntity func() botsfw.BotAppUser) GaeAppUserStore {
//	return GaeAppUserStore{
//		appUserEntityType: appUserEntityType,
//		appUserEntityKind: appUserEntityKind,
//		newUserEntity:     newUserEntity,
//		GaeBaseStore:      NewGaeBaseStore(appUserEntityKind),
//	}
//}
//
//// ************************** Helper functions **************************
//
//func (s GaeAppUserStore) appUserKey(c context.Context, appUserID int64) *datastore.Key {
//	return datastore.NewKey(c, s.appUserEntityKind, "", appUserID, nil)
//}
//
//// ************************** Implementations of  botsfw.AppUserStore **************************
//
//// GetAppUserByID returns application user ID
//func (s GaeAppUserStore) GetAppUserByID(c context.Context, appUserID int64, appUser botsfw.BotAppUser) error {
//	if appUserID == 0 {
//		panic("appUserID == 0")
//	}
//	return nds.Get(c, s.appUserKey(c, appUserID), appUser)
//}
//
//// CreateAppUser creates app user entity in DB
//func (s GaeAppUserStore) CreateAppUser(c context.Context, botID string, actor botsfw.WebhookActor) (int64, botsfw.BotAppUser, error) {
//	return s.createAppUser(c, botID, actor)
//}
//
//func (s GaeAppUserStore) createAppUser(c context.Context, botID string, actor botsfw.WebhookActor) (int64, botsfw.BotAppUser, error) {
//	appUserEntity := s.newUserEntity()
//	appUserEntity.SetBotUserID(actor.Platform(), botID, fmt.Sprintf("%v", actor.GetID()))
//	appUserEntity.SetNames(actor.GetFirstName(), actor.GetLastName(), actor.GetUserName())
//	key, err := nds.Put(c, s.appUserKey(c, 0), appUserEntity)
//	return key.IntID(), appUserEntity, err
//}
//
//func (s GaeAppUserStore) getAppUserIDByBotUserKey(c context.Context, botUserKey *datastore.Key) (int64, error) {
//	query := datastore.NewQuery(s.appUserEntityKind).Filter("TelegramUserIDs =", botUserKey.IntID()).KeysOnly().Limit(2)
//	//appUsers := reflect.MakeSlice(reflect.SliceOf(s.appUserEntityType), 0, 2)
//	keys, err := query.GetAll(c, nil)
//	if err != nil {
//		log.Errorf(c, "Failed to query app users by TelegramUserIDs: %v", err)
//		return 0, err
//	}
//	switch len(keys) {
//	case 0:
//		return 0, nil
//	case 1:
//		return keys[0].IntID(), nil
//	default:
//		return 0, fmt.Errorf("Found few app users by %v", botUserKey)
//	}
//}
//
////func (s GaeAppUserStore) SaveAppUser(c context.Context, appUserID int64, appUserEntity botsfw.BotAppUser) error {
////	if appUserID == 0 {
////		panic("appUserID == 0")
////	}
////	_, err := nds.Put(c, s.appUserKey(c, appUserID), appUserEntity)
////	return err
////}
