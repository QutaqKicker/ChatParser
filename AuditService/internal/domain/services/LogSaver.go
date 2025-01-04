package services

import (
	"audit/internal/domain/models"
	"database/sql"
	"github.com/QutaqKicker/ChatParser/common/dbHelper"
	auditv1 "github.com/QutaqKicker/ChatParser/protos/gen/go/audit"
	"log"
	"time"
)

type LogSaver struct {
	log *log.Logger
	db  *sql.DB
}

func NewLogSaver(log *log.Logger, db *sql.DB) *LogSaver {
	return &LogSaver{log: log, db: db}
}

var insertStatement *sql.Stmt

func (s *LogSaver) SaveLog(serviceName string, auditType auditv1.AuditType, message string) {
	if insertStatement == nil {
		var err error
		insertStatement, err = s.db.Prepare(dbHelper.BuildInsert[models.Log](false))

		if err != nil {
			log.Fatal(err)
		}
	}

	_, err := insertStatement.Exec(serviceName, auditType, message, time.Now())
	if err != nil {
		log.Fatal(err)
	}
}
