package actions

import (
	"Client/internal/utils"
	"github.com/QutaqKicker/ChatParser/Router/pkg/routerClient"
)

var MainActions = []utils.Action[func(routerClient *routerClient.RouterClient)]{
	{"Выйти", nil},
	{"Найти сообщения", GetMessages},
	{"Получить количество сообщений по пользователям", GetUsersWithMessagesCount},
	{"Импортировать сообщения", ImportMessages},
}
