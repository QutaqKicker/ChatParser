package main

import (
	"Client/internal/actions"
	"Client/internal/utils"
)

func main() {
	//routerClient := struct{}{} //TODO set router client

	for {
		mainActionCallback := utils.ShowActionsForSelect(actions.MainActions)
		if mainActionCallback == nil { //Пришла команда на выход
			break
		}

		mainActionCallback()
	}
}
