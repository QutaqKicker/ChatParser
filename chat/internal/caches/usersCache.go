package caches

import (
	"chat/internal/db"
	"chat/internal/domain/queryFilters"
	"database/sql"
)

var UsersCache = newUsersCache()

type usersCache CacheOfNames[string]

func cacheInitializer(tx *sql.Tx) {
	rows, err := tx.Query(db.BuildQuery[queryFilters.UserFilter](db.QueryBuildRequest[queryFilters.UserFilter]{}))
	if err != nil {
		panic(err)
	}

	for (row := rows.)
}

func newUsersCache() *usersCache {
	return &usersCache{}
}
