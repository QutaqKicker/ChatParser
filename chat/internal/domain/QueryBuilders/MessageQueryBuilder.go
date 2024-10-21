package QueryBuilders

import (
	"chat/internal/domain/filters"
	"chat/internal/domain/models"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func BuildQuery(filter *filters.MessageFilter) string {
	queryBuilder := strings.Builder{}

	if filter.SpecifySelect != "" {
		queryBuilder.WriteString(fmt.Sprintf("select %s", filter.SpecifySelect))
	} else {
		queryBuilder.WriteString(BuildSelect[models.Message]())
	}

	queryBuilder.WriteString(BuildWhere(filter))
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

type MessageUpdateValue struct {
	Field          reflect.StructField
	NewValueInt    int
	NewValueString string
	NewValueTime   time.Time
}

func BuildUpdate(filter *filters.MessageFilter, values []MessageUpdateValue) string {
	updateBuilder := strings.Builder{}
	updateBuilder.WriteString("update messages")
	for i, value := range values {
		if i > 0 {
			updateBuilder.WriteString(",")
		}
		updateBuilder.WriteString("\n set ")
		sqlColumn := value.Field.Tag.Get("column")

		if sqlColumn != "" {
			updateBuilder.WriteString(fmt.Sprintf("%s = ", sqlColumn))
		} else {
			updateBuilder.WriteString(fmt.Sprintf("%s = ", value.Field.Name))
		}

		switch value.Field.Type.Name() {
		case "int32":
			updateBuilder.WriteString(strconv.Itoa(value.NewValueInt))
		case "string":
			updateBuilder.WriteString(value.NewValueString)
		case "time.Time": //TODO ??
			updateBuilder.WriteString(value.NewValueTime.Format("YYYY.MM.DD hh.mm.ss"))
		}
	}

	updateBuilder.WriteString(BuildWhere(filter))

	return updateBuilder.String()
}

func BuildDelete(filter *filters.MessageFilter) string {
	deleteBuilder := strings.Builder{}
	deleteBuilder.WriteString("delete from messages")
	deleteBuilder.WriteString(BuildWhere(filter))
	return deleteBuilder.String()
}

func BuildWhere(filter *filters.MessageFilter) string {
	whereBuilder := strings.Builder{}
	whereBuilder.WriteString("\n where 1 == 1")

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

func BuildSorter(filter *filters.MessageFilter) string {
	if filter.Sorts != nil {
		return fmt.Sprintf("\n order by %s", strings.Join(filter.Sorts, ", "))
	}
	return "\n order by created desc"
}

func RowToEntity() models.Message {
	//TODO Разобраться что приходит из БД и сделать парс
	return models.Message{}
}
