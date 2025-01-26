package services

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/QutaqKicker/ChatParser/Common/dbHelper"
	"log/slog"
	"math"
	"time"
	"users/internal/domain/filters"
	"users/internal/domain/models"
)

type UserMessageCounter struct {
	log *slog.Logger
	db  *sql.DB
}

func NewUserMessageCounter(log *slog.Logger, db *sql.DB) *UserMessageCounter {
	return &UserMessageCounter{log: log, db: db}
}

func (c *UserMessageCounter) UpdateUserMessagesCount(ctx context.Context, userName string, count int) {
	tx, err := c.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		c.log.Error(err.Error())
		return
	}

	userFilter := filters.NewUserFilter().WhereName(userName)

	selectQuery, selectParams := dbHelper.BuildQuery[models.User](dbHelper.SelectBuildRequest{
		Filter:        userFilter,
		SelectType:    dbHelper.Special,
		SpecialSelect: "messages_count",
	})

	rows, err := tx.QueryContext(ctx, selectQuery, selectParams...)
	if err != nil {
		c.log.Error(err.Error())
		return
	}

	oldMessageCount := math.MinInt32
	for rows.Next() {
		err := rows.Scan(&oldMessageCount)
		if err != nil {
			fmt.Println(err)
		}
	}

	if err != nil {
		c.log.Error(err.Error())
		return
	}

	if oldMessageCount == math.MinInt32 {
		insertQuery := dbHelper.BuildInsert[models.User](false)
		newUser := &models.User{Name: userName, MessagesCount: count, Created: time.Now()}
		_, err := tx.ExecContext(ctx, insertQuery, newUser.FieldValuesAsArray()...)
		if err != nil {
			c.log.Error(err.Error())
		}
	} else {
		//TODO Надо бы перепилить апдейтер dbHelper чтобы можно было апдейтить поля на основании имеющегося значения, потом можно будет избавиться от транзакции
		updateUserQuery, updateUserParams := dbHelper.BuildUpdate[models.User](
			dbHelper.SetUpdate("messages_count", oldMessageCount+count),
			userFilter)

		_, err = c.db.ExecContext(ctx, updateUserQuery, updateUserParams...)
		if err != nil {
			c.log.Error(err.Error())
		}
	}
	tx.Commit()
}
