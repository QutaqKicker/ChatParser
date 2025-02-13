package main

import (
	"github.com/QutaqKicker/ChatParser/Common/constants"
	"github.com/QutaqKicker/ChatParser/Common/myKafka"
	"github.com/QutaqKicker/ChatParser/Common/myLogs"
	"os"
	"os/signal"
	"router/internal/router"
	"strconv"
	"syscall"
)

func main() {

	auditProducer := myKafka.NewAuditProducer()
	defer auditProducer.Close()

	logger := myLogs.SetupLogger(auditProducer, "Router")

	logger.Info("started router")

	routerPort, _ := strconv.Atoi(os.Getenv(constants.RouterPortEnvName))

	chatPort, _ := strconv.Atoi(os.Getenv(constants.ChatPortEnvName))
	userPort, _ := strconv.Atoi(os.Getenv(constants.UserPortEnvName))
	backupPort, _ := strconv.Atoi(os.Getenv(constants.BackupPortEnvName))

	router := router.NewRouter(logger)

	go func() {
		err := application.Run()
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	application.Stop()
}
