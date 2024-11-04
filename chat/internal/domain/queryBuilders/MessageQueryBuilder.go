package queryBuilders

import (
	"chat/internal/domain/filters"
	"chat/internal/domain/models"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Entity interface {
	TableName() string
}

func BuildQuery(request *filters.QueryBuildRequest) string {
	queryBuilder := strings.Builder{}

	queryBuilder.WriteString(BuildSelect[*models.Message](request.SelectType, request.SpecialSelect))

	queryBuilder.WriteString(BuildWhere(request.Filter))
	queryBuilder.WriteString(BuildSorter(request.Sorter))
	return queryBuilder.String()
}

func BuildSelect[T Entity](selectType filters.SelectType, specialSelect string) string {
	if specialSelect == "" {
		specialSelect = "*"
	}

	switch selectType {
	case filters.All:
		return fmt.Sprintf("select %s", ColumnNamesWithAliases[T]())
	case filters.Sum:
		return fmt.Sprintf("select sum(%s)", specialSelect)
	case filters.Count:
		return fmt.Sprintf("select count(%s)", specialSelect)
	default:
		return fmt.Sprintf("select %s", specialSelect)
	}
}

func ColumnNames[T Entity]() string {
	t := reflect.TypeOf(*new(T))
	sqlColumns := make([]string, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.Tag.Get("not-mapped") == "true" {
			continue
		}

		sqlColumn := field.Tag.Get("column")
		if sqlColumn == "" {
			sqlColumn = field.Name
		}

		sqlColumns = append(sqlColumns, sqlColumn)
	}

	return strings.Join(sqlColumns, ", ")
}

func ColumnNamesWithAliases[T Entity]() string {
	t := reflect.TypeOf(new(T))

	sqlColumns := make([]string, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.Tag.Get("not-mapped") == "true" {
			continue
		}

		sqlColumn := field.Tag.Get("column")
		if sqlColumn != "" {
			sqlColumns = append(sqlColumns, fmt.Sprintf("%s as %s", sqlColumn, field.Name))
		} else {
			sqlColumns = append(sqlColumns, field.Name)
		}
	}

	return strings.Join(sqlColumns, ", ")
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

	if !filter.MinCreatedDate.IsZero() {
		whereBuilder.WriteString(fmt.Sprintf("\n and %s < created", filter.MinCreatedDate.Format("YYYY.MM.DD")))
	}

	if !filter.MaxCreatedDate.IsZero() {
		whereBuilder.WriteString(fmt.Sprintf("\n and created < %s", filter.MinCreatedDate.Format("YYYY.MM.DD")))
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

func BuildInsert[T Entity](withReturning bool) string {
	insertQuery := strings.Builder{}
	t := *new(T)
	insertQuery.WriteString(fmt.Sprintf("insert into %s", T.TableName(t)))
	insertQuery.WriteString(fmt.Sprintf("\n\t(%s)", ColumnNames[T]()))

	entityType := reflect.TypeOf(t)

	values := make([]string, 0, entityType.NumField())
	currParamIndex := 1
	for i := 0; i < entityType.NumField(); i++ {
		if entityType.Field(i).Tag.Get("not-mapped") == "true" {
			continue
		}

		values = append(values, fmt.Sprintf("$%d", currParamIndex))
		currParamIndex++
	}

	insertQuery.WriteString(fmt.Sprintf("\nvalues\n\t(%s)", strings.Join(values, ", ")))

	if withReturning {
		insertQuery.WriteString(fmt.Sprintf("\nreturning %s", ColumnNamesWithAliases[T]()))
	}

	return insertQuery.String()
}

func BuildSorter(sorter []string) string {
	if sorter != nil {
		return fmt.Sprintf("\n order by %s", strings.Join(sorter, ", "))
	}
	return "\n order by created desc"
}

func RowToEntity() models.Message {
	//TODO Разобраться что приходит из БД и сделать парс
	return models.Message{}
}
