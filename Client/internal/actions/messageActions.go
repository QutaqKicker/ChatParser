package actions

import (
	"Client/internal/utils"
	"context"
	"fmt"
	chatv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/chat"
	commonv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/common"
	"github.com/QutaqKicker/ChatParser/Router/pkg/routerClient"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
	"time"
)

func GetMessages(router *routerClient.RouterClient) {
	filter := new(commonv1.MessagesFilter)
	for {
		filterActionCallback := utils.ShowActionsForSelect(FilterActions)
		if filterActionCallback == nil { //Пришла команда на выход
			break
		}

		err := filterActionCallback(filter)
		if err != nil {
			fmt.Println(err)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	count, err := router.GetMessagesCount(ctx, filter)
	if err != nil {
		panic(err)
	}

	if count == 0 {
		fmt.Println("Сообщений не обнаружено")
		return
	} else {
		fmt.Printf("Обнаружено %d сообщений. Что с ними сделать?\n", count)

		for {
			messageActionCallback := utils.ShowActionsForSelect(MessageActions)
			if messageActionCallback == nil { //Пришла команда на выход
				break
			}

			ok := messageActionCallback(ctx, router, filter)
			if ok {
				break
			}
		}
	}
}

var FilterActions = []utils.Action[func(*commonv1.MessagesFilter) error]{
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

var MessageActions = []utils.Action[func(ctx context.Context, router *routerClient.RouterClient, f *commonv1.MessagesFilter) bool]{
	{"Ничего", nil},
	{"Вывести в консоль", func(ctx context.Context, router *routerClient.RouterClient, f *commonv1.MessagesFilter) bool {
		messagesResponse, err := router.SearchMessages(ctx, &chatv1.SearchMessagesRequest{Filter: f})
		if err != nil {
			fmt.Printf("Повторите попытку. Возникла ошибка: %v", err)
			return false
		} else {
			for i, message := range messagesResponse.Messages {
				fmt.Printf("%d - %s. От %s: %s", i, message.Created, message.UserName, message.Text)
			}
			return true
		}
	}},
	{"Экспортировать", func(ctx context.Context, router *routerClient.RouterClient, f *commonv1.MessagesFilter) bool {
		fmt.Println("Экспортируем в формате csv(1) или parquet(2)? Введите 1 или 2")
		var exportType string
		_, err := fmt.Scanln(&exportType)
		if err != nil {
			fmt.Printf("Некорректный ввод, повторите попытку. Ошибка: %v", err)
			return false
		} else {
			if trimmedExportType := strings.TrimSpace(exportType); trimmedExportType == "1" {
				exportType = "csv"
			} else if trimmedExportType == "2" {
				exportType = "parquet"
			} else {
				fmt.Printf("Некорректный ввод, повторите попытку. Ошибка: %v", err)
				return false
			}
		}
		ok, err := router.ExportToDir(ctx, f, exportType)
		if err != nil {
			fmt.Printf("Не получилось. Повторите попытку. Ошибка: %v\n", err)
			return ok
		}
		return ok
	}},
}
