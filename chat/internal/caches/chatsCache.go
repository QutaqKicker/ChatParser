package caches

import (
	"chat/internal/dbHelper"
	"chat/internal/domain/filters"
	"chat/internal/domain/models"
	"sync"
	"time"
)

// ChatsCache Кэш чатов. Ключ - имя чата, значение - ключ чата
var ChatsCache = newChatsCache()

type chatsCache struct {
	CacheOfNames[int32]
}

func newChatsCache() *chatsCache {
	return &chatsCache{
		CacheOfNames[int32]{
			mutex:       sync.RWMutex{},
			once:        sync.Once{},
			initializer: chatsCacheInitializer,
			dbUpdater:   chatsCacheDbUpdater,
			dbInserter:  chatsCacheDbInserter,
		},
	}
}

func chatsCacheInitializer(querier dbOrTx, elems *map[string]int32) {
	query, _ := dbHelper.BuildQuery[models.Chat](dbHelper.QueryBuildRequest{})
	rows, err := querier.Query(query)
	if err != nil {
		panic(err)
	}

	chats, err := dbHelper.RowsToEntities[models.Chat](rows)
	if err != nil {
		panic(err)
	}

	alreadyExistsChats := make(map[string]int32, len(chats))
	for _, user := range chats {
		alreadyExistsChats[user.Name] = user.Id
	}
	*elems = alreadyExistsChats
}

func chatsCacheDbUpdater(tx dbOrTx, oldKey int32, newKey int32) {
	updateUserQuery, userParams := dbHelper.BuildUpdate[models.Chat](dbHelper.SetUpdate("id", newKey),
		filters.NewChatFilter().WhereId(oldKey))
	tx.Exec(updateUserQuery, userParams...)

	updateMessagesQuery, messageParams := dbHelper.BuildUpdate[models.Chat](dbHelper.SetUpdate("user_id", newKey),
		filters.NewMessageFilter().WhereChatIds([]int32{oldKey}))

	tx.Exec(updateMessagesQuery, messageParams...)
}

func chatsCacheDbInserter(tx dbOrTx, name string, key int32) int32 {
	if key == 0 {
		newChat := models.Chat{Name: name, Created: time.Now()}
		insertQuery := dbHelper.BuildInsert[models.User](false, true)
		rows, _ := tx.Query(insertQuery, newChat.Name, newChat.Created)
		for rows.Next() {
			rows.Scan(&key)
		}
		return key
	} else {
		newChat := models.Chat{Id: key, Name: name, Created: time.Now()}
		insertQuery := dbHelper.BuildInsert[models.User](true, false)
		tx.Exec(insertQuery, newChat.FieldValuesAsArray())
		return key
	}
}
