package services

import (
	"audit/internal/domain/models"
	"database/sql"
	"github.com/QutaqKicker/ChatParser/common/dbHelper"
	"log"
	"log/slog"
	"time"
)

type LogSaver struct {
	log *slog.Logger
	db  *sql.DB
}

func NewLogSaver(log *slog.Logger, db *sql.DB) *LogSaver {
	return &LogSaver{log: log, db: db}
}

var insertStatement *sql.Stmt

func (s *LogSaver) SaveLog(serviceName string, auditType int, message string) {
	if insertStatement == nil {
		var err error
		insertStatement, err = s.db.Prepare(dbHelper.BuildInsert[models.Log](true))

		if err != nil {
			log.Fatal(err)
		}
	}

	_, err := insertStatement.Exec(serviceName, auditType, message, time.Now())
	if err != nil {
		log.Fatal(err)
	}
}
