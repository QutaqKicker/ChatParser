package caches

import "database/sql"

var UsersCache = newUsersCache()

type usersCache CacheOfNames[string]

func cacheInitializer(tx *sql.Tx) {
	stmt, err := tx.Prepare()
}

func newUsersCache() *usersCache {
	return &usersCache{}
}
