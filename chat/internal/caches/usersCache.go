package caches

import (
	"chat/internal/db"
	"chat/internal/domain/models"
	"sync"
)

var UsersCache = newUsersCache()

type usersCache CacheOfNames[string]

func usersCacheInitializer(querier dbOrTx) {
	rows, err := querier.Query(db.BuildQuery[models.User](db.QueryBuildRequest{}))
	if err != nil {
		panic(err)
	}

	users, err := db.RowsToEntities[models.User](rows)
	if err != nil {
		panic(err)
	}

	UsersCache.elems = make(map[string]string, len(users))
	for _, user := range users {
		UsersCache.elems[user.Name] = user.Id
	}
}

func usersCacheDbUpdater(tx dbOrTx, oldKey string, newKey string) {
	//tx.Exec(db.BuildUpdate())
}

func usersCacheDbInserter(tx dbOrTx, name string, key string) {

}

func newUsersCache() *usersCache {
	return &usersCache{
		mutex:       sync.RWMutex{},
		once:        sync.Once{},
		initializer: usersCacheInitializer,
		dbUpdater:   usersCacheDbUpdater,
		dbInserter:  usersCacheDbInserter,
	}
}
