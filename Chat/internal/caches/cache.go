package caches

import (
	"database/sql"
	"sync"
)

// CacheOfNames Кэш, для поиска ключей по имени сущностей. Необходим, т.к. в html формате не хранятся id пользователей
// и чатов, а в других форматах хранятся. Когда появляется конкретный id чата или пользователя, его нужно актуализировать
// в базе и в уже импортированных сообщениях с этим юзером
type CacheOfNames[T comparable] struct {
	elems map[string]T
	mutex sync.RWMutex
	once  sync.Once
	// initializer единожды запускается при взятии или внесении. Он загружает в кэш все имеющиеся сущности из БД
	initializer func(tx dbOrTx, elems *map[string]T) //TODO Не уверен что будет пахать. Надо проверить
	// dbUpdater обновляет в БД старый ключ на новый, если в выгрузке пришел конкретный айди сущности
	dbUpdater func(tx dbOrTx, oldKey T, newKey T)
	// dbInserter вставляет в БД новый ключ и название сущности, раз сущности с таким названием в БД не существует
	dbInserter func(tx dbOrTx, name string, key T) T
}

type dbOrTx interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...any) (sql.Result, error)
}

// GetByName Получить ключ сущности по его имени
func (c *CacheOfNames[T]) GetByName(tx dbOrTx, name string) (key T, ok bool) {
	c.mutex.RLock()
	c.once.Do(func() { c.initializer(tx, &c.elems) })
	key, ok = c.elems[name]
	c.mutex.RUnlock()
	return
}

// GetByKey Получить имя сущности по его ключу
func (c *CacheOfNames[T]) GetByKey(tx dbOrTx, key T) (name string, ok bool) {
	c.mutex.RLock()
	c.once.Do(func() { c.initializer(tx, &c.elems) })
	for i, value := range c.elems {
		if value == key {
			name = i
		}
	}
	ok = name != ""

	c.mutex.RUnlock()
	return
}

// Set Внести в кэш указанное имя и айди сущности и апсертнуть эту сущность в базе
func (c *CacheOfNames[T]) Set(tx dbOrTx, name string, key T) T {
	c.mutex.Lock()
	c.once.Do(func() { c.initializer(tx, &c.elems) })
	if oldKey, ok := c.elems[name]; ok {
		if oldKey != key {
			c.dbUpdater(tx, oldKey, key)
		}
	} else {
		key = c.dbInserter(tx, name, key)
	}
	c.elems[name] = key
	c.mutex.Unlock()

	return key
}
