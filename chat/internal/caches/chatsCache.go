package caches

import (
	"chat/internal/domain/filters"
	"chat/internal/domain/models"
	"github.com/QutaqKicker/ChatParser/dbHelper/pkg"
	"log"
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
	query, _ := pkg.BuildQuery[models.Chat](pkg.SelectBuildRequest{})
	rows, err := querier.Query(query)
	if err != nil {
		panic(err)
	}

	chats, err := pkg.RowsToEntities[models.Chat](rows)
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
	updateUserQuery, userParams := pkg.BuildUpdate[models.Chat](pkg.SetUpdate("id", newKey),
		filters.NewChatFilter().WhereId(oldKey))
	tx.Exec(updateUserQuery, userParams...)

	updateMessagesQuery, messageParams := pkg.BuildUpdate[models.Chat](pkg.SetUpdate("user_id", newKey),
		filters.NewMessageFilter().WhereChatIds([]int32{oldKey}))

	tx.Exec(updateMessagesQuery, messageParams...)
}

func chatsCacheDbInserter(tx dbOrTx, name string, key int32) int32 {
	if key == 0 {
		newChat := models.Chat{Name: name, Created: time.Now()}
		insertQuery := pkg.BuildInsert[models.Chat](true)
		rows, err := tx.Query(insertQuery, newChat.Name, newChat.Created)
		if err != nil {
			panic(err)
		}
		for rows.Next() {
			err = rows.Scan(&key)
			if err != nil {
				log.Fatal(err)
			}
		}
		return key
	} else {
		newChat := models.Chat{Id: key, Name: name, Created: time.Now()}
		insertQuery := pkg.BuildInsert[models.Chat](false)
		tx.Exec(insertQuery, newChat.FieldValuesAsArray())
		return key
	}
}
