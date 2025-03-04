package caches

import (
	"chat/internal/domain/filters"
	"chat/internal/domain/models"
	"fmt"
	"github.com/QutaqKicker/ChatParser/Common/dbHelper"
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
	query, _ := dbHelper.BuildQuery[models.User](dbHelper.SelectBuildRequest{})
	rows, err := querier.Query(query)
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
	_, err := tx.Exec(updateUserQuery, userParams...)
	if err != nil {
		fmt.Println(err)
	}
}

func usersCacheDbInserter(tx dbOrTx, name string, key string) string {
	if key == "" {
		newUser := models.User{Id: name, Name: name, Created: time.Now()}
		insertQuery := dbHelper.BuildInsert[models.User](false)
		_, err := tx.Exec(insertQuery, newUser.FieldValuesAsArray()...)
		if err != nil {
			fmt.Println(err)
		}

		return newUser.Id
	} else {
		newUser := models.User{Id: key, Name: name, Created: time.Now()}
		insertQuery := dbHelper.BuildInsert[models.User](false)
		_, err := tx.Exec(insertQuery, newUser.FieldValuesAsArray()...)
		if err != nil {
			fmt.Println(err)
		}

		return key
	}
}
