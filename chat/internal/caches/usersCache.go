package caches

import (
	"chat/internal/domain/filters"
	"chat/internal/domain/models"
	"github.com/QutaqKicker/ChatParser/dbHelper/pkg"
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
	query, _ := pkg.BuildQuery[models.User](pkg.SelectBuildRequest{})
	rows, err := querier.Query(query)
	if err != nil {
		panic(err)
	}

	users, err := pkg.RowsToEntities[models.User](rows)
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
	updateUserQuery, userParams := pkg.BuildUpdate[models.User](pkg.SetUpdate("id", newKey),
		filters.NewUserFilter().WhereId(oldKey))
	tx.Exec(updateUserQuery, userParams...)

	updateMessagesQuery, messageParams := pkg.BuildUpdate[models.User](pkg.SetUpdate("user_id", newKey),
		filters.NewMessageFilter().WhereUserIds([]string{oldKey}))

	tx.Exec(updateMessagesQuery, messageParams...)
}

func usersCacheDbInserter(tx dbOrTx, name string, key string) string {
	if key == "" {
		newUser := models.User{Id: name, Name: name, Created: time.Now()}
		insertQuery := pkg.BuildInsert[models.User](false)
		_, err := tx.Exec(insertQuery, newUser.FieldValuesAsArray()...)
		if err != nil {
			panic(err)
		}

		return newUser.Id
	} else {
		newUser := models.User{Id: key, Name: name, Created: time.Now()}
		insertQuery := pkg.BuildInsert[models.User](false)
		tx.Exec(insertQuery, newUser.FieldValuesAsArray()...)
		return key
	}
}
