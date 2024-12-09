package caches

import (
	"database/sql"
	"sync"
)

type CacheOfNames[T comparable] struct {
	elems       map[string]T
	mutex       sync.RWMutex
	once        sync.Once
	initializer func(dbOrTx)
	dbUpdater   func(tx dbOrTx, oldKey T, newKey T)
	dbInserter  func(tx dbOrTx, name string, key T) T
}

type dbOrTx interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...any) (sql.Result, error)
}

func (c *CacheOfNames[T]) GetKeyByName(tx dbOrTx, name string) (key T, ok bool) {
	c.mutex.RLock()
	c.once.Do(func() { c.initializer(tx) })
	key, ok = c.elems[name]
	c.mutex.RUnlock()
	return
}

func (c *CacheOfNames[T]) SetNewChat(tx dbOrTx, name string, key T) {
	c.mutex.Lock()
	c.once.Do(func() { c.initializer(tx) })
	if oldKey, ok := c.elems[name]; ok {
		if oldKey != key {
			c.dbUpdater(tx, oldKey, key)
		}
	} else {
		key = c.dbInserter(tx, name, key)
	}
	c.elems[name] = key
	c.mutex.Unlock()
}
