package caches

import (
	"database/sql"
	"sync"
)

type CacheOfNames[T comparable] struct {
	elems       map[string]T
	mutex       sync.RWMutex
	once        sync.Once
	initializer func(*sql.Tx)
	dbUpdater   func(tx *sql.Tx, oldKey T, newKey T)
	dbInserter  func(tx *sql.Tx, name string, key T)
}

func (c *CacheOfNames[T]) Get(tx *sql.Tx, name string) (key T, ok bool) {
	c.mutex.RLock()
	c.once.Do(func() { c.initializer(tx) })
	key, ok = c.elems[name]
	c.mutex.RUnlock()
	return
}

func (c *CacheOfNames[T]) Set(tx *sql.Tx, name string, key T) {
	c.mutex.Lock()
	c.once.Do(func() { c.initializer(tx) })
	if oldKey, ok := c.elems[name]; ok {
		if oldKey != key {
			c.dbUpdater(tx, oldKey, key)
		}
	} else {
		c.dbInserter(tx, name, key)
	}
	c.elems[name] = key
	c.mutex.Unlock()
}
