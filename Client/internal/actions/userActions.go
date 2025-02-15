package actions

import (
	"context"
	"fmt"
	"github.com/QutaqKicker/ChatParser/Router/pkg/routerClient"
)

func GetUsersWithMessagesCount(router *routerClient.RouterClient) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	usersResponse, err := router.GetUsersWithMessagesCount(ctx)
	if err != nil {
		fmt.Printf("Не получилось. Повторите попытку. Ошибка: %v", err)
	}

	for _, user := range usersResponse.Users {
		fmt.Printf("Пользователь '%s' понаписал %d сообщений\n", user.Name, user.MessagesCount)
	}
}
