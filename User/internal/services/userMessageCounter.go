package services

import (
	"context"
	"database/sql"
	"github.com/QutaqKicker/ChatParser/Common/dbHelper"
	"log"
	"time"
	"users/internal/domain/filters"
	"users/internal/domain/models"
)

type UserMessageCounter struct {
	log *log.Logger
	db  *sql.DB
}

func NewUserMessageCounter(log *log.Logger, db *sql.DB) *UserMessageCounter {
	return &UserMessageCounter{log: log, db: db}
}

func (c *UserMessageCounter) UpdateUserMessagesCount(ctx context.Context, userName string, count int) {
	tx, err := c.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		c.log.Fatal(err)
		return
	}

	userFilter := filters.NewUserFilter().WhereName(userName)

	selectQuery, selectParams := dbHelper.BuildQuery[models.User](dbHelper.SelectBuildRequest{
		Filter:        userFilter,
		SelectType:    dbHelper.Special,
		SpecialSelect: "messages_count",
	})

	rows, err := tx.QueryContext(ctx, selectQuery, selectParams)
	if err != nil {
		c.log.Fatal(err)
		return
	}

	if rows == nil {
		insertQuery := dbHelper.BuildInsert[models.User](false)
		newUser := &models.User{Name: userName, MessagesCount: count, Created: time.Now()}
		_, err := tx.ExecContext(ctx, insertQuery, newUser.FieldValuesAsArray())
		if err != nil {
			c.log.Fatal(err)
		}
		return
	}

	users, err := dbHelper.RowsToEntities[models.User](rows)
	if err != nil {
		c.log.Fatal(err)
		return
	}

	currentMessageCount := users[0].MessagesCount

	//TODO Надо бы перепилить апдейтер dbHelper чтобы можно было апдейтить поля на основании имеющегося значения, потом можно будет избавиться от транзакции
	updateUserQuery, updateUserParams := dbHelper.BuildUpdate[models.User](
		dbHelper.SetUpdate("messages_count", currentMessageCount+count),
		userFilter)

	_, err = c.db.ExecContext(ctx, updateUserQuery, updateUserParams)
	if err != nil {
		c.log.Fatal(err)
	}
}
