package db

import (
	"chat/internal/domain/queryFilters"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type Entity interface {
	TableName() string
}

func BuildQuery[T Entity](request QueryBuildRequest) (string, []interface{}) {
	queryBuilder := strings.Builder{}

	queryBuilder.WriteString(buildSelect[T](request.SelectType, request.SpecialSelect))

	whereString, values := buildWhere(request.Filter, 0)
	queryBuilder.WriteString(whereString)
	queryBuilder.WriteString(buildSorter(request.Sorter))
	return queryBuilder.String(), values
}

type UpdateValue struct {
	FieldName string
	NewValue  interface{}
}

// TODO Сделать удобнее, или вообще переместить в сущности отдельными методами
func BuildUpdate(values []UpdateValue, filter any) (string, []interface{}) {
	updateBuilder := strings.Builder{}
	updateBuilder.WriteString("update messages")
	queryValues := make([]interface{}, 0)
	paramIndex := 1
	updateBuilder.WriteString("\n set ")
	sets := make([]string, 0)

	for _, value := range values {
		sets = append(sets, fmt.Sprintf("%s = $%d", value.FieldName, paramIndex))
		paramIndex++
		queryValues = append(queryValues, value.NewValue)
	}

	updateBuilder.WriteString(strings.Join(sets, ",\n"))

	whereString, whereValues := buildWhere(filter, paramIndex)
	queryValues = append(queryValues, whereValues...)
	updateBuilder.WriteString(whereString)

	return updateBuilder.String(), queryValues
}

func BuildDelete[T Entity](filter *queryFilters.MessageFilter) (string, []interface{}) {
	deleteBuilder := strings.Builder{}
	deleteBuilder.WriteString(fmt.Sprintf("delete from %s", (*new(T)).TableName()))
	whereString, values := buildWhere(filter, 0)
	deleteBuilder.WriteString(whereString)
	return deleteBuilder.String(), values
}

func BuildInsert[T Entity](withReturning bool) string {
	insertQuery := strings.Builder{}
	t := *new(T)
	insertQuery.WriteString(fmt.Sprintf("insert into %s", T.TableName(t)))
	insertQuery.WriteString(fmt.Sprintf("\n\t(%s)", columnNames[T]()))

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
		insertQuery.WriteString(fmt.Sprintf("\nreturning %s", columnNamesWithAliases[T]()))
	}

	return insertQuery.String()
}

func RowsToEntities[T Entity](rows *sql.Rows) ([]T, error) {
	entities := make([]T, 0)
	entityType := reflect.TypeOf(new(T))
	mappedFieldsIndexes := make([]int, 0)

	for i := 0; i < entityType.NumField(); i++ {
		if entityType.Field(i).Tag.Get("not-mapped") == "true" {
			continue
		} else {
			mappedFieldsIndexes = append(mappedFieldsIndexes, i)
		}
	}

	for rows != nil {
		entity := new(T)
		entityValue := reflect.ValueOf(&entity).Elem()

		columns := make([]interface{}, 0, len(mappedFieldsIndexes))
		for _, index := range mappedFieldsIndexes {
			field := entityValue.Field(index)
			columns = append(columns, field.Addr().Interface())
		}

		err := rows.Scan(columns...)

		if err != nil {
			return nil, err
		}

		entities = append(entities, *entity)
	}
	return entities, nil
}

func buildSelect[T Entity](selectType SelectType, specialSelect string) string {
	if specialSelect == "" {
		specialSelect = "*"
	}

	switch selectType {
	case All:
		return fmt.Sprintf("select %s", columnNamesWithAliases[T]())
	case Sum:
		return fmt.Sprintf("select sum(%s)", specialSelect)
	case Count:
		return fmt.Sprintf("select count(%s)", specialSelect)
	default:
		return fmt.Sprintf("select %s", specialSelect)
	}
}

func columnNames[T Entity]() string {
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

func columnNamesWithAliases[T Entity]() string {
	t := reflect.TypeOf(*new(T))

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

func buildWhere[TFilter any](filter TFilter, firstParamIndex int) (string, []interface{}) {
	whereBuilder := strings.Builder{}
	whereBuilder.WriteString("\n where 1 == 1")

	values := make([]interface{}, 0)
	t := reflect.TypeOf(filter)
	v := reflect.ValueOf(filter)
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			fieldValue := v.Field(i)
			if fieldValue.IsZero() {
				continue
			}
			columnName := t.Field(i).Tag.Get("column")

			relation := t.Field(i).Tag.Get("relation")
			switch relation {
			case "=", "<", ">":
				whereBuilder.WriteString(fmt.Sprintf("\n  and %s %s $%d", columnName, relation, len(values)+firstParamIndex))
				values = append(values, fieldValue.Interface())
			case "like":
				whereBuilder.WriteString(fmt.Sprintf("\n  and %s like '%%$%d%%'", columnName, len(values)+firstParamIndex))
				values = append(values, fieldValue.String())
			case "in":
				inParams := make([]string, 0, fieldValue.Len())
				switch inSlice := fieldValue.Interface().(type) {
				case []string:
					for _, value := range inSlice {
						inParams = append(inParams, fmt.Sprintf("$%d", len(values)+firstParamIndex))
						values = append(values, value)
					}
				case []int:
					for _, value := range inSlice {
						inParams = append(inParams, fmt.Sprintf("$%d", len(values)+firstParamIndex))
						values = append(values, value)
					}
				default:
					panic("unknown type")
				}
				whereBuilder.WriteString(fmt.Sprintf("\n  and %s in (%s)", columnName, strings.Join(inParams, ", ")))
			}
		}
	}

	return whereBuilder.String(), values
}

func buildSorter(sorter []string) string {
	if sorter != nil {
		return fmt.Sprintf("\n order by %s", strings.Join(sorter, ", "))
	}
	return "\n order by created desc"
}
