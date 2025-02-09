package actions

import (
	"Client/internal/utils"
	"fmt"
	commonv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/common"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
	"time"
)

var GetMessagesActions = []utils.Action[func(*commonv1.MessagesFilter) error]{
	{"Фильтр готов", nil},
	{"Указать Id", func(f *commonv1.MessagesFilter) error {
		id, err := utils.ScanInt()
		if err != nil {
			return err
		}
		f.Id = int32(id)
		return nil
	}},
	{"Указать период создания", func(f *commonv1.MessagesFilter) error {
		dates := []struct {
			message    string
			filterDate **timestamppb.Timestamp
		}{
			{"Введите минимальную дату создания в формате yyyy-mm-dd", &f.MinCreatedDate},
			{"Введите максимальную дату создания в формате yyyy-mm-dd", &f.MaxCreatedDate},
		}
		for _, value := range dates {
			fmt.Println(value.message)

			var rawDate string
			_, err := fmt.Scanln(&rawDate)
			if err != nil {
				return err
			}

			date, err := time.Parse("2006-01-02", rawDate)
			if err != nil {
				return err
			}
			*value.filterDate = timestamppb.New(date)
		}

		return nil
	}},
	{"Указать искомый фрагмент текста", func(f *commonv1.MessagesFilter) error {
		var rawSubtext string
		_, err := fmt.Scanln(&rawSubtext)
		if err != nil {
			return err
		}

		f.SubText = strings.TrimSpace(rawSubtext)
		return nil
	}},
}

func GetMessages() {
	filter := commonv1.MessagesFilter{}
	for {
		getMessagesActionCallback := utils.ShowActionsForSelect(GetMessagesActions)
		if getMessagesActionCallback == nil { //Пришла команда на выход
			break
		}

		err := getMessagesActionCallback(&filter)
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println(filter)

	//TODO send that filter to chatservice

	//TODO show count of messageActions via that filter and ask user, what he need: show, export or delete that messageActions
}
