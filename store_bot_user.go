package gae

//
//import (
//	"context"
//	"fmt"
//	"github.com/bots-go-framework/bots-fw/botsfw"
//	"github.com/strongo/log"
//	"github.com/strongo/nds"
//	"google.golang.org/appengine/v2/datastore"
//	"time"
//)
//
//// GaeBotUserStore persist user to GAE datastore
//type GaeBotUserStore struct {
//	GaeBaseStore
//	botUserKey                func(c context.Context, botUserID interface{}) *datastore.Key
//	validateBotUserEntityType func(entity botsfw.BotUser)
//	newBotUserEntity          func(apiUser botsfw.WebhookActor) botsfw.BotUser
//	gaeAppUserStore           GaeAppUserStore
//}
//
//var _ botsfw.BotUserStore = (*GaeBotUserStore)(nil) // Check for interface implementation at compile time
//
//// ************************** Implementations of  botsfw.BotUserStore **************************
//
//// GetBotUserByID returns bot user by ID
//func (s GaeBotUserStore) GetBotUserByID(c context.Context, botUserID interface{}) (botsfw.BotUser, error) { // Former LoadBotUserEntity
//	botUserEntity := s.newBotUserEntity(nil)
//	err := nds.Get(c, s.botUserKey(c, botUserID), botUserEntity)
//	if err == datastore.ErrNoSuchEntity {
//		return nil, nil
//	}
//	return botUserEntity, err
//}
//
//// SaveBotUser saves bot user by ID
//func (s GaeBotUserStore) SaveBotUser(c context.Context, botUserID interface{}, userEntity botsfw.BotUser) error { // Former SaveBotUserEntity
//	// TODO: Architecture needs refactoring as it not transactional save
//	// We load bot user entity outside of here (out of transaction) and save here. It can change since then.
//	s.validateBotUserEntityType(userEntity)
//	userEntity.SetUpdatedTime(time.Now())
//	err := nds.RunInTransaction(c, func(c context.Context) error {
//		key := s.botUserKey(c, botUserID)
//		existingBotUser := s.newBotUserEntity(nil)
//		err := nds.Get(c, key, existingBotUser)
//		if err != nil {
//			if err == datastore.ErrNoSuchEntity {
//				err = nil
//			}
//			return err
//		} else {
//			if existingBotUser.GetAppUserIntID() != userEntity.GetAppUserIntID() {
//				return fmt.Errorf(
//					"Data integrity issue, existingBotUser.GetAppUserIntID():%v != userEntity.GetAppUserIntID():%v",
//					existingBotUser.GetAppUserIntID(),
//					userEntity.GetAppUserIntID(),
//				)
//			}
//		}
//		_, err = nds.Put(c, key, userEntity)
//		if err != nil {
//			err = fmt.Errorf("failed to put user entity to datastore in SaveBotUser(): %w", err)
//		}
//		return err
//	}, nil)
//	return err
//}
//
//// CreateBotUser creates bot user
//func (s GaeBotUserStore) CreateBotUser(c context.Context, botID string, apiUser botsfw.WebhookActor) (botsfw.BotUser, error) {
//	log.Debugf(c, "GaeBotUserStore.CreateBotUser(botID=%v, apiUser=%T) started...", botID, apiUser)
//	botUserID := apiUser.GetID()
//	botUserEntity := s.newBotUserEntity(apiUser)
//
//	var (
//		appUserID    int64
//		isNewAppUser bool
//	)
//
//	err := nds.RunInTransaction(c, func(ctx context.Context) (err error) {
//		botUserKey := s.botUserKey(ctx, botUserID)
//
//		switch err = nds.Get(ctx, botUserKey, botUserEntity); err {
//		case nil:
//			appUserID = botUserEntity.GetAppUserIntID()
//			return
//		case datastore.ErrNoSuchEntity:
//			// We need to create entity, so continue execution
//		default: // err != nil && err != datastore.ErrNoSuchEntity
//			return
//		}
//
//		// First try to search for existing users by Bot ID in case if entity was deleted
//		if appUserID, err = s.gaeAppUserStore.getAppUserIDByBotUserKey(c, botUserKey); err != nil {
//			return
//		}
//
//		if appUserID == 0 { // This is most expected case
//			appUserID, _, err = s.gaeAppUserStore.createAppUser(ctx, botID, apiUser)
//			if err != nil {
//				log.Errorf(c, "Failed to create app user: %v", err)
//				return
//			}
//			if appUserID == 0 {
//				panic("appUserID == 0")
//			}
//			isNewAppUser = true
//		}
//
//		botUserEntity.SetAppUserIntID(appUserID)
//		botUserEntity.SetUpdatedTime(time.Now())
//
//		if _, err = nds.Put(ctx, botUserKey, botUserEntity); err != nil {
//			return
//		}
//		return
//	}, &datastore.TransactionOptions{XG: true})
//
//	if err != nil {
//		return nil, fmt.Errorf("failed to create bot user: %w", err)
//	}
//
//	log.Debugf(c, "GaeBotUserStore.CreateBotUser() => appUserID: %v, isNewAppUser: %v", appUserID, isNewAppUser)
//
//	return botUserEntity, nil
//}
