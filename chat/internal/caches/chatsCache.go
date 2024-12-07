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
	CacheOfNames[int]
}

func newChatsCache() *chatsCache {
	return &chatsCache{
		CacheOfNames[int]{
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

	ChatsCache.elems = make(map[string]int, len(chats))
	for _, user := range chats {
		ChatsCache.elems[user.Name] = user.Id
	}
}

func ChatsCacheDbUpdater(tx dbOrTx, oldKey int, newKey int) {
	updateUserQuery, userParams := dbHelper.BuildUpdate[models.Chat](dbHelper.SetUpdate("id", newKey),
		filters.NewChatFilter().WhereId(oldKey))
	tx.Exec(updateUserQuery, userParams...)

	updateMessagesQuery, messageParams := dbHelper.BuildUpdate[models.Chat](dbHelper.SetUpdate("user_id", newKey),
		filters.NewMessageFilter().WhereChatIds([]int{oldKey}))

	tx.Exec(updateMessagesQuery, messageParams...)
}

func ChatsCacheDbInserter(tx dbOrTx, name string, key int) {
	newUser := models.Chat{Id: key, Name: name, Created: time.Now()}
	insertQuery := dbHelper.BuildInsert[models.User](false)
	tx.Exec(insertQuery, newUser.FieldValuesAsArray())

}
