package caches

import (
	"chat/internal/dbHelper"
	"chat/internal/domain/filters"
	"chat/internal/domain/models"
	"sync"
	"time"
)

var ChatsCache = newChatsCache()

type chatsCache struct {
	CacheOfNames[int32]
}

func newChatsCache() *chatsCache {
	return &chatsCache{
		CacheOfNames[int32]{
			mutex:       sync.RWMutex{},
			once:        sync.Once{},
			initializer: ChatsCacheInitializer,
			dbUpdater:   ChatsCacheDbUpdater,
			dbInserter:  ChatsCacheDbInserter,
		},
	}
}

func ChatsCacheInitializer(querier dbOrTx) {
	rows, err := querier.Query(dbHelper.BuildQuery[models.Chat](dbHelper.QueryBuildRequest{}))
	if err != nil {
		panic(err)
	}

	chats, err := dbHelper.RowsToEntities[models.Chat](rows)
	if err != nil {
		panic(err)
	}

	ChatsCache.elems = make(map[string]int32, len(chats))
	for _, user := range chats {
		ChatsCache.elems[user.Name] = user.Id
	}
}

func ChatsCacheDbUpdater(tx dbOrTx, oldKey int32, newKey int32) {
	updateUserQuery, userParams := dbHelper.BuildUpdate[models.Chat](dbHelper.SetUpdate("id", newKey),
		filters.NewChatFilter().WhereId(oldKey))
	tx.Exec(updateUserQuery, userParams...)

	updateMessagesQuery, messageParams := dbHelper.BuildUpdate[models.Chat](dbHelper.SetUpdate("user_id", newKey),
		filters.NewMessageFilter().WhereChatIds([]int32{oldKey}))

	tx.Exec(updateMessagesQuery, messageParams...)
}

func ChatsCacheDbInserter(tx dbOrTx, name string, key int32) int32 {
	if key == 0 {
		newChat := models.Chat{Name: name, Created: time.Now()}
		insertQuery := dbHelper.BuildInsert[models.User](false, true)
		rows, _ := tx.Query(insertQuery, newChat.Name, newChat.Created)
		rows.Scan(&key)
		return key
	} else {
		newChat := models.Chat{Id: key, Name: name, Created: time.Now()}
		insertQuery := dbHelper.BuildInsert[models.User](true, false)
		tx.Exec(insertQuery, newChat.FieldValuesAsArray())
		return key
	}
}
