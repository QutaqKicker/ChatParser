package actions

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

func ImportMessages() {
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

		//TODO sending to chatservice dirPath
	}

}
