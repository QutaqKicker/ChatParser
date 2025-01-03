package services

import (
	"context"
	"database/sql"
	"github.com/QutaqKicker/ChatParser/common/dbHelper"
	"log"
	"users/internal/domain/filters"
	"users/internal/domain/models"
)

type userMessageCounter struct {
	log *log.Logger
	db  *sql.DB
}

func (c *userMessageCounter) UpdateUserMessagesCount(ctx context.Context, userId, userName string, count int) {
	//TODO Надо бы перепилить апдейтер dbHelper чтобы можно было апдейтить поля на основании имеющегося значения
	updateUserQuery, userParams := dbHelper.BuildUpdate[models.User](dbHelper.SetUpdate("id", newKey),
		filters.NewUserFilter().WhereId(oldKey))

	query, queryParams := dbHelper.BuildUpdate[models.User](dbHelper.SetUpdate())
	c.db.ExecContext(ctx)
}
