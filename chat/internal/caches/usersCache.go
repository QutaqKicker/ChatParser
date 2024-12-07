package caches

import (
	"chat/internal/dbHelper"
	"chat/internal/domain/filters"
	"chat/internal/domain/models"
	"sync"
	"time"
)

var UsersCache = newUsersCache()

type usersCache struct {
	CacheOfNames[string]
}

func newUsersCache() *usersCache {
	return &usersCache{
		CacheOfNames[string]{
			mutex:       sync.RWMutex{},
			once:        sync.Once{},
			initializer: usersCacheInitializer,
			dbUpdater:   usersCacheDbUpdater,
			dbInserter:  usersCacheDbInserter,
		},
	}
}

func usersCacheInitializer(querier dbOrTx) {
	rows, err := querier.Query(dbHelper.BuildQuery[models.User](dbHelper.QueryBuildRequest{}))
	if err != nil {
		panic(err)
	}

	users, err := dbHelper.RowsToEntities[models.User](rows)
	if err != nil {
		panic(err)
	}

	UsersCache.elems = make(map[string]string, len(users))
	for _, user := range users {
		UsersCache.elems[user.Name] = user.Id
	}
}

func usersCacheDbUpdater(tx dbOrTx, oldKey string, newKey string) {
	updateUserQuery, userParams := dbHelper.BuildUpdate[models.User](dbHelper.SetUpdate("id", newKey),
		filters.NewUserFilter().WhereId(oldKey))
	tx.Exec(updateUserQuery, userParams...)

	updateMessagesQuery, messageParams := dbHelper.BuildUpdate[models.User](dbHelper.SetUpdate("user_id", newKey),
		filters.NewMessageFilter().WhereUserIds([]string{oldKey}))

	tx.Exec(updateMessagesQuery, messageParams...)
}

func usersCacheDbInserter(tx dbOrTx, name string, key string) {
	newUser := models.User{Id: key, Name: name, Created: time.Now()}
	insertQuery := dbHelper.BuildInsert[models.User](false)
	tx.Exec(insertQuery, newUser.FieldValuesAsArray())
}
