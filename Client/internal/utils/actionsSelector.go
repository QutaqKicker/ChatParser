package utils

import (
	"fmt"
)

type Action[TCallback any] struct {
	Name     string
	Callback TCallback
}

func ShowActionsForSelect[TCallback any](actions []Action[TCallback]) TCallback {
	fmt.Println("Введите цифру требуемого действия:")

	for i, action := range actions {
		fmt.Printf("%d - %s\n", i, action.Name)
	}

	for {
		selectedIndex, err := ScanInt()
		if err != nil || selectedIndex < 0 || selectedIndex >= len(actions) {
			fmt.Println("Указанное значение невалидно или отсутствует в перечне доступных действий. Повторите попытку.")
		} else {
			return actions[selectedIndex].Callback
		}
	}
}
