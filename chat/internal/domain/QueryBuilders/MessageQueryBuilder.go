package QueryBuilders

import (
	"chat/internal/domain/models"
	"fmt"
	"reflect"
	"strings"
)

func BuildQuery(filter *models.MessageFilter) string {
	queryBuilder := strings.Builder{}

	if filter.SpecifySelect != "" {
		queryBuilder.WriteString(fmt.Sprintf("select %s", filter.SpecifySelect))
	} else {
		queryBuilder.WriteString(BuildSelect[models.Message]())
	}

	queryBuilder.WriteString("\n where 1 == 1")
	queryBuilder.WriteString(BuildFilter(filter))
	queryBuilder.WriteString(BuildSorter(filter))
	return queryBuilder.String()
}

func BuildSelect[T any]() string {
	selectQuery := strings.Builder{}
	selectQuery.WriteString("select ")
	t := reflect.TypeOf(new(T))
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		sqlColumn := field.Tag.Get("column")

		if sqlColumn != "" {
			selectQuery.WriteString(fmt.Sprintf("%s as %s", sqlColumn, field.Name))
		} else {
			selectQuery.WriteString(field.Name)
		}
	}

	return selectQuery.String()
}

func BuildFilter(filter *models.MessageFilter) string {
	whereBuilder := strings.Builder{}

	if !filter.MinDate.IsZero() {
		whereBuilder.WriteString(fmt.Sprintf("\n and %s < created", filter.MinDate.Format("YYYY.MM.DD")))
	}

	if !filter.MaxDate.IsZero() {
		whereBuilder.WriteString(fmt.Sprintf("\n and created < %s", filter.MinDate.Format("YYYY.MM.DD")))
	}

	if filter.SubText != "" {
		whereBuilder.WriteString(fmt.Sprintf("\n and text like '%%%s%%'", filter.SubText))
	}

	if len(filter.UserIds) > 0 {
		userIdsFilter := strings.Builder{}
		for _, value := range filter.UserIds {
			userIdsFilter.WriteString(fmt.Sprintf("%d ,", value))
		}
		whereBuilder.WriteString(fmt.Sprintf("\n and user_id in (%s)", userIdsFilter))
	}

	if len(filter.ChatIds) > 0 {
		chatIdsFilter := strings.Builder{}
		for _, value := range filter.ChatIds {
			chatIdsFilter.WriteString(fmt.Sprintf("%d ,", value))
		}
		whereBuilder.WriteString(fmt.Sprintf("\n and chat_id in (%s)", chatIdsFilter))
	}

	return whereBuilder.String()
}

func BuildSorter(filter *models.MessageFilter) string {
	if filter.Sorts != nil {
		return fmt.Sprintf("\n order by %s", strings.Join(filter.Sorts, ", "))
	}
	return "\n order by created desc"
}

func RowToEntity() models.Message {
	//TODO Разобраться что приходит из БД и сделать парс
}
