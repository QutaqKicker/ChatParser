package postgres

import (
	"chat/internal/domain/queryFilters"
	"chat/internal/domain/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("postgres", "psqlInfo")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveMessages(ctx context.Context, []models.Message) error {
	const op = "storage.postgres.SaveMessages"

	stmt, err := s.db.Prepare("INSERT ATYATYA VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, " ", " ")
	if err != nil {
		var sqlErr pgx.PgError

		if errors.As(err, &sqlErr) && sqlErr.Code == pgx.ErrNoRows.Error() {
			return errors.New("test")
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) Messages(ctx context.Context, filter queryFilters.QueryBuildRequest) ([]models.Message, error) {
	const op = "storage.postgres.Messages"

	stmt, err := s.db.Prepare("select where ????")
	if err != nil{

	}

	rows, err := stmt.QueryContext(ctx, " ")

	for _, value := range rows.{
		var message models.Message

		value
	}
}
