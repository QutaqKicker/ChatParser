package caches

import (
	"chat/internal/dbHelper"
	"chat/internal/domain/filters"
	"chat/internal/domain/models"
	"sync"
	"time"
)

// UsersCache Кэш юзеров. Ключ - имя юзера, значение - ключ юзера
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

func usersCacheInitializer(querier dbOrTx, elems *map[string]string) {
	rows, err := querier.Query(dbHelper.BuildQuery[models.User](dbHelper.QueryBuildRequest{}))
	if err != nil {
		panic(err)
	}

	users, err := dbHelper.RowsToEntities[models.User](rows)
	if err != nil {
		panic(err)
	}

	alreadyExistsUsers := make(map[string]string, len(users))
	for _, user := range users {
		alreadyExistsUsers[user.Name] = user.Id
	}
	*elems = alreadyExistsUsers
}

func usersCacheDbUpdater(tx dbOrTx, oldKey string, newKey string) {
	updateUserQuery, userParams := dbHelper.BuildUpdate[models.User](dbHelper.SetUpdate("id", newKey),
		filters.NewUserFilter().WhereId(oldKey))
	tx.Exec(updateUserQuery, userParams...)

	updateMessagesQuery, messageParams := dbHelper.BuildUpdate[models.User](dbHelper.SetUpdate("user_id", newKey),
		filters.NewMessageFilter().WhereUserIds([]string{oldKey}))

	tx.Exec(updateMessagesQuery, messageParams...)
}

func usersCacheDbInserter(tx dbOrTx, name string, key string) string {
	if key == "" {
		newUser := models.User{Name: name, Created: time.Now()}
		insertQuery := dbHelper.BuildInsert[models.User](false, true)
		rows, _ := tx.Query(insertQuery, newUser.Name, newUser.Created)
		rows.Scan(&key)
		return key
	} else {
		newUser := models.User{Id: key, Name: name, Created: time.Now()}
		insertQuery := dbHelper.BuildInsert[models.User](true, false)
		tx.Exec(insertQuery, newUser.FieldValuesAsArray())
		return key
	}
}
