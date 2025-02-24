package main

import (
	"Client/internal/actions"
	"Client/internal/utils"
	"fmt"
	"github.com/QutaqKicker/ChatParser/Common/constants"
	"github.com/QutaqKicker/ChatParser/Router/pkg/routerClient"
	"os"
)

func main() {
	routerPort := os.Getenv(constants.RouterPortEnvName)
	router := routerClient.NewRouterClient(fmt.Sprintf("http://localhost:%s", routerPort))

	for {
		mainActionCallback := utils.ShowActionsForSelect(actions.MainActions)
		if mainActionCallback == nil { //Пришла команда на выход
			break
		}
		mainActionCallback(router)
	}
}
