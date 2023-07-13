package gae

//
//import (
//	"context"
//	"fmt"
//	"github.com/bots-go-framework/bots-fw/botsfw"
//	"github.com/strongo/log"
//	"github.com/strongo/nds"
//	"google.golang.org/appengine/v2/datastore"
//	"strconv"
//	"time"
//)
//
//// EntityTypeValidator validate entity types
//type EntityTypeValidator interface {
//}
//
//// GaeBotChatStore persists chat to GAE datastore
//type GaeBotChatStore struct {
//	GaeBaseStore
//	botChats                  map[string]botsfw.BotChat
//	NewBotChatKey             func(c context.Context, botID, botChatId string) *datastore.Key
//	validateBotChatEntityType func(entity botsfw.BotChat)
//	newBotChatEntity          func() botsfw.BotChat
//}
//
//func NewGaeBotChatStore(
//	entityKind string,
//	newBotChatKey func(c context.Context, botID, botChatId string) *datastore.Key,
//	validateBotChatEntityType func(entity botsfw.BotChat),
//	newBotChatEntity func() botsfw.BotChat,
//) *GaeBotChatStore {
//	if newBotChatKey == nil {
//		newBotChatKey = func(c context.Context, botID, botChatId string) *datastore.Key {
//			return datastore.NewKey(c, entityKind, botsfw.NewChatID(botID, botChatId), 0, nil)
//		}
//	}
//	return &GaeBotChatStore{
//		GaeBaseStore:              NewGaeBaseStore(entityKind),
//		NewBotChatKey:             newBotChatKey,
//		newBotChatEntity:          newBotChatEntity,
//		validateBotChatEntityType: validateBotChatEntityType,
//	}
//}
//
//var _ botsfw.BotChatStore = (*GaeBotChatStore)(nil) // Check for interface implementation at compile time
//
//// ************************** Implementations of  botsfw.ChatStore **************************
//
//// GetBotChatEntityByID returns bot chat entity by ID
//func (s *GaeBotChatStore) GetBotChatEntityByID(c context.Context, botID, botChatID string) (botChatEntity botsfw.BotChat, err error) { // Former LoadBotChatEntity
//	//log.Debugf(c, "GaeBotChatStore.GetBotChatEntityByID(%v)", botChatId)
//	if s.botChats == nil {
//		s.botChats = make(map[string]botsfw.BotChat, 1)
//	}
//	botChatEntity = s.newBotChatEntity()
//	botChatKey := s.NewBotChatKey(c, botID, botChatID)
//	//c, _ = context.WithDeadline(c, time.Now().Add(time.Second))
//	if err = nds.Get(c, botChatKey, botChatEntity); err != nil {
//		if err != datastore.ErrNoSuchEntity {
//			return
//		}
//		if s.entityKind == "TgChat" { // TODO: Remove workaround to fix old entities
//			var tgChatID int64
//			if tgChatID, err = strconv.ParseInt(botChatID, 10, 64); err != nil {
//				err = fmt.Errorf("failed to parse botChatID to int: %w", err)
//				return
//			}
//			intKey := datastore.NewKey(c, s.entityKind, "", tgChatID, nil)
//			if err = nds.Get(c, intKey, botChatEntity); err != nil {
//				if err == datastore.ErrNoSuchEntity {
//					log.Infof(c, fmt.Errorf("there is no bot chat entity with intID=%d: %w", intKey.IntID(), err).Error())
//					err = botsfw.ErrEntityNotFound
//				}
//				return
//			}
//			log.Infof(c, "Telegram chat entity Found by int ID, will attempt to migrate...")
//			if err = nds.RunInTransaction(c, func(tc context.Context) (err error) {
//				if err = nds.Get(tc, intKey, botChatEntity); err != nil {
//					return
//				}
//				if err = nds.Delete(tc, intKey); err != nil {
//					return
//				}
//				if _, err = nds.Put(tc, botChatKey, botChatEntity); err != nil {
//					return
//				}
//				return
//			}, &datastore.TransactionOptions{XG: true}); err != nil {
//				log.Errorf(c, fmt.Errorf("failed to migrate Telegram chat entity: %w", err).Error())
//				return
//
//			}
//			log.Infof(c, "Telegram chat entity migrated to new key: [%v]", botChatKey.StringID())
//		}
//		err = botsfw.ErrEntityNotFound
//		return
//	}
//	if err == nil {
//		s.botChats[botChatKey.StringID()] = botChatEntity
//	}
//	return
//}
//
//// SaveBotChat saves bot chat
//func (s *GaeBotChatStore) SaveBotChat(c context.Context, botID, botChatID string, chatEntity botsfw.BotChat) error { // Former SaveBotChatEntity
//	s.validateBotChatEntityType(chatEntity)
//	chatEntity.SetUpdatedTime(time.Now())
//	_, err := nds.Put(c, s.NewBotChatKey(c, botID, botChatID), chatEntity)
//	return err
//}
//
//// NewBotChatEntity creates new bot chat entity
//func (s *GaeBotChatStore) NewBotChatEntity(c context.Context, botID string, botChat botsfw.WebhookChat, appUserID string, botUserID string, isAccessGranted bool) botsfw.BotChat {
//	botChatID := botChat.GetID()
//	log.Debugf(c, "NewBotChatEntity(botID=%v, botChatID=%v, appUserID=%v, botUserID=%v, isAccessGranted=%v)", botID, botChatID, appUserID, botUserID, isAccessGranted)
//	botChatEntity := s.newBotChatEntity()
//	botChatEntity.SetBotID(botID)
//
//	if botChat.IsGroupChat() {
//		botChatEntity.SetIsGroupChat(true)
//	} else {
//		if appUserIntID, err := strconv.ParseInt(appUserID, 10, 64); err == nil { // TODO: change to STRING APP USER ID
//			panic(fmt.Sprintf("appUserID is not valid int64: %s", appUserID))
//		} else {
//			botChatEntity.SetAppUserIntID(appUserIntID)
//		}
//		botChatEntity.SetBotUserID(botUserID)
//	}
//
//	botChatEntity.SetAccessGranted(isAccessGranted)
//	s.botChats[s.NewBotChatKey(c, botID, botChatID).StringID()] = botChatEntity // TODO: No need to create a key instance, create dedicated func to create ID?
//	return botChatEntity
//}
//
//// Close is called on request completion
//func (s *GaeBotChatStore) Close(c context.Context) error { // Former SaveBotChatEntity
//	if len(s.botChats) == 0 {
//		log.Debugf(c, "GaeBotChatStore.Close(): Nothing to save")
//		return nil
//	}
//	//log.Debugf(c, "GaeBotChatStore.Close(): %v entities to save", len(s.botChats))
//	var chatKeys []*datastore.Key
//	var chatEntities []botsfw.BotChat
//	now := time.Now()
//	for chatID, chatEntity := range s.botChats {
//		s.validateBotChatEntityType(chatEntity)
//		chatEntity.SetUpdatedTime(now)
//		chatEntity.SetDtLastInteraction(now)
//		chatKeys = append(chatKeys, datastore.NewKey(c, s.entityKind, chatID, 0, nil))
//		chatEntities = append(chatEntities, chatEntity)
//	}
//	_, err := nds.PutMulti(c, chatKeys, chatEntities)
//	if err == nil {
//		//log.Debugf(c, "Successfully saved %v BotChat entities with keys: %v", len(chatKeys), chatKeys)
//		s.botChats = nil
//	} else {
//		log.Errorf(c, "Failed to save %v BotChat entities: %v", len(chatKeys), err)
//	}
//	return err
//}
