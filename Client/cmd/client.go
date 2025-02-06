package main

import (
	"fmt"
	commonv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/common"
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	for {
		actionCallback := showActionsForSelect()
		actionCallback()
	}
}

func showActionsForSelect() func() {
	fmt.Println("Введите цифру требуемого действия:")

	for i, action := range actions {
		fmt.Printf("%d - %s\n", i, action.name)
	}

	for {
		selectedIndex, err := scanInt()
		if err != nil || selectedIndex < 0 || selectedIndex >= len(actions) {
			fmt.Println("Указанное значение невалидно или отсутствует в перечне доступных действий. Повторите попытку.")
		} else {
			return actions[selectedIndex].callback
		}
	}
}

func scanInt() (int, error) {
	var scanned string
	_, err := fmt.Scanln(&scanned)
	if err != nil {
		return 0, err
	}

	result, err := strconv.Atoi(strings.TrimSpace(scanned))
	if err != nil {
		return 0, err
	}

	return result, nil
}

var actions = []struct {
	name     string
	callback func()
}{
	{"Выйти", Exit},
	{"Найти сообщения", GetMessages},
	{"Получить количество сообщений по пользователям", GetUsersWithMessagesCount},
	{"Экспортировать сообщения", ExportMessages},
	{"Импортировать сообщения", ImportMessages},
}

func Exit() {
	os.Exit(0)
}

func GetMessages() {
	filter := commonv1.MessagesFilter{}
	for {
		getMessagesActions
	}
}

var getMessagesActions = []struct {
	name   string
	action func(*commonv1.MessagesFilter, any) error
}{
	{"Указать Id", func(f *commonv1.MessagesFilter, value any) error {
		id, err := scanInt()
		if err != nil {
			return err
		}
		f.Id = int32(id)
		return nil
	}},
	{"Указать период создания", func(f *commonv1.MessagesFilter, value any) error {
		dates := []struct {
			message    string
			filterDate *timestamppb.Timestamp
		}{
			{"Введите минимальную дату создания в формате yyyy-mm-dd", f.MinCreatedDate},
			{"Введите максимальную дату создания в формате yyyy-mm-dd", f.MaxCreatedDate},
		}
		for _, value := range dates {
			fmt.Println(value.message)

			var rawDate string
			_, err := fmt.Scanln(rawDate)
			if err != nil {
				return err
			}

			date, err := time.Parse("2006-01-02", rawDate)
			if err != nil {
				return err
			}
			value.filterDate = timestamppb.New(date)
		}

		return nil
	}},
	{"Указать искомый подтекст", func(f *commonv1.MessagesFilter, value any) error {
		var rawSubtext string
		_, err := fmt.Scanln(&rawSubtext)
		if err != nil {
			return err
		}

		f.SubText = strings.TrimSpace(rawSubtext)
		return nil
	}},
}

func GetUsersWithMessagesCount() {

}

func ExportMessages() {

}

func ImportMessages() {

}
