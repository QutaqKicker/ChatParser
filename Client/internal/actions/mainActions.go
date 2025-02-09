package actions

import (
	"Client/internal/utils"
)

var MainActions = []utils.Action[func()]{
	{"Выйти", nil},
	{"Найти сообщения", GetMessages},
	{"Получить количество сообщений по пользователям", GetUsersWithMessagesCount},
	{"Импортировать сообщения", ImportMessages},
}
