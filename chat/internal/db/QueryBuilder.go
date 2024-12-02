package db

import (
	"chat/internal/domain/models"
	"chat/internal/domain/queryFilters"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Entity interface {
	TableName() string
}

func BuildQuery[TFilter any](request QueryBuildRequest[TFilter]) (string, []interface{}) {
	queryBuilder := strings.Builder{}

	queryBuilder.WriteString(BuildSelect[models.Message](request.SelectType, request.SpecialSelect))

	whereString, values := BuildWhere(request.Filter)
	queryBuilder.WriteString(whereString)
	queryBuilder.WriteString(BuildSorter(request.Sorter))
	return queryBuilder.String(), values
}

func BuildSelect[T Entity](selectType SelectType, specialSelect string) string {
	if specialSelect == "" {
		specialSelect = "*"
	}

	switch selectType {
	case All:
		return fmt.Sprintf("select %s", ColumnNamesWithAliases[T]())
	case Sum:
		return fmt.Sprintf("select sum(%s)", specialSelect)
	case Count:
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

type UpdateValue struct {
	Field          reflect.StructField
	NewValueInt    int
	NewValueString string
	NewValueTime   time.Time
}

func BuildUpdate(filter *queryFilters.MessageFilter, values []UpdateValue) (string, []interface{}) {
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

	whereString, queryValues := BuildWhere(filter)
	updateBuilder.WriteString(whereString)

	return updateBuilder.String(), queryValues
}

func BuildDelete[T Entity](filter *queryFilters.MessageFilter) (string, []interface{}) {
	deleteBuilder := strings.Builder{}
	deleteBuilder.WriteString(fmt.Sprintf("delete from %s", (*new(T)).TableName()))
	whereString, values := BuildWhere(filter)
	deleteBuilder.WriteString(whereString)
	return deleteBuilder.String(), values
}

func BuildWhere[TFilter any](filter TFilter) (string, []interface{}) {
	whereBuilder := strings.Builder{}
	whereBuilder.WriteString("\n where 1 == 1")

	values := make([]interface{}, 0)
	t := reflect.TypeOf(filter)
	v := reflect.ValueOf(filter)
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			fieldValue := v.Field(i)
			columnName := t.Field(i).Tag.Get("column")

			relation := t.Field(i).Tag.Get("relation")
			switch relation {
			case "=":
			case "<":
			case ">":
				values = append(values, fieldValue.Interface())
				whereBuilder.WriteString(fmt.Sprintf("\n  and %s %s $%d", columnName, relation, len(values)))
			case "like":
				values = append(values, fieldValue.String())
				whereBuilder.WriteString(fmt.Sprintf("\n  and %s like '%%$%d%%'", columnName, len(values)))
			case "in":
				inParams := make([]string, 0, fieldValue.Len())
				switch inSlice := fieldValue.Interface().(type) {
				case []string:
					for _, value := range inSlice {
						values = append(values, value)
						inParams = append(inParams, fmt.Sprintf("$%d", len(values)))
					}
				case []int:
					for _, value := range inSlice {
						values = append(values, value)
						inParams = append(inParams, fmt.Sprintf("$%d", len(values)))
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

/*package main

import (
    "fmt"
    "reflect"
)

type Address struct {
    City  string
    State string
}

type User struct {
    Name    string
    Age     int
    Address // Встраиваем Address
}

func main() {
    user := User{
        Name: "Alice",
        Age:  30,
        Address: Address{
            City:  "Wonderland",
            State: "Fantasy",
        },
    }

    // Используем рефлексию для перебора полей
    val := reflect.ValueOf(user)
    typ := reflect.TypeOf(user)

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)
        fieldType := typ.Field(i)

        // Если поле встроенной структуры, выведем его с правильным именем
        if fieldType.Anonymous {
            for j := 0; j < field.NumField(); j++ {
                includedField := field.Field(j)
                includedFieldType := field.Type().Field(j)

                fmt.Printf("%s: %v\n", includedFieldType.Name, includedField.Interface())
            }
        } else {
            fmt.Printf("%s: %v\n", fieldType.Name, field.Interface())
        }
    }
}*/

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
