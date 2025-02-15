package main

import (
	"Router/internal/router"
	"fmt"
	"github.com/QutaqKicker/ChatParser/Common/constants"
	"github.com/QutaqKicker/ChatParser/Common/myKafka"
	"github.com/QutaqKicker/ChatParser/Common/myLogs"
	"net/http"
	"os"
)

func main() {

	auditProducer := myKafka.NewAuditProducer()
	defer auditProducer.Close()

	logger := myLogs.SetupLogger(auditProducer, "Router")

	logger.Info("started router")

	routerPort := os.Getenv(constants.RouterPortEnvName)

	chatPort := os.Getenv(constants.ChatPortEnvName)
	userPort := os.Getenv(constants.UserPortEnvName)
	backupPort := os.Getenv(constants.BackupPortEnvName)

	routerMux := router.NewRouter(logger, chatPort, userPort, backupPort)

	err := http.ListenAndServe(fmt.Sprintf(":%s", routerPort), routerMux)
	if err != nil {
		logger.Error(err.Error())
	}
}
