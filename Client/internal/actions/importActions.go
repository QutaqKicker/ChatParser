package actions

import (
	"context"
	"errors"
	"fmt"
	"github.com/QutaqKicker/ChatParser/Router/pkg/routerClient"
	"io/fs"
	"os"
)

func ImportMessages(router *routerClient.RouterClient) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("Введите путь к файлам, которые нужно сымпортировать")
	for {
		var importPath string
		_, err := fmt.Scan(&importPath)
		if err != nil {
			fmt.Printf("Не получилось. Ошибка: %v. Повторите попытку \n", err)
			continue
		}

		_, err = os.Stat(importPath)
		if err != nil && errors.Is(err, fs.ErrNotExist) {
			fmt.Println("Такой директории не существует. Повторите попытку")
			continue
		}

		ok, err := router.ParseFromDir(ctx, importPath)
		if !ok || err != nil {
			fmt.Printf("Не получилось. Ошибка: %v. Повторите попытку \n", err)
			continue
		} else {
			break
		}
	}
}
