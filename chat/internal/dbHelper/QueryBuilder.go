package dbHelper

import (
	"chat/internal/domain/filters"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type Entity interface {
	TableName() string
}

// BuildQuery Построить селект для типа из дженерика согласно настройкам запроса в QueryBuildRequest
func BuildQuery[T Entity](request QueryBuildRequest) (string, []interface{}) {
	queryBuilder := strings.Builder{}

	queryBuilder.WriteString(buildSelect[T](request.SelectType, request.SpecialSelect))

	t := *new(T)
	queryBuilder.WriteString(fmt.Sprintf("\nfrom %s", T.TableName(t)))

	whereString, values := buildWhere(request.Filter, 0)
	queryBuilder.WriteString(whereString)
	//queryBuilder.WriteString(buildSorter(request.SortColumnIndexes))//TODO
	if request.Take != 0 {
		queryBuilder.WriteString(fmt.Sprintf("\n LIMIT %d", request.Take))
	}
	if request.Skip != 0 {
		queryBuilder.WriteString(fmt.Sprintf("\n OFFSET %d", request.Skip))
	}
	//queryBuilder.WriteString()
	return queryBuilder.String(), values
}

type UpdateValue struct {
	FieldName string
	NewValue  any
}

type UpdateValues []UpdateValue

func SetUpdate(fieldName string, newValue any) UpdateValues {
	return []UpdateValue{{fieldName, newValue}}
}

func (u UpdateValues) AndUpdate(fieldName string, newValue any) UpdateValues {
	u = append(u, UpdateValue{fieldName, newValue})
	return u
}

// TODO Сделать удобнее, или вообще переместить в сущности отдельными методами
func BuildUpdate[T Entity](values UpdateValues, filter any) (string, []interface{}) {
	updateBuilder := strings.Builder{}
	t := *new(T)
	updateBuilder.WriteString(fmt.Sprintf("update %s", T.TableName(t)))
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

func BuildDelete[T Entity](filter *filters.MessageFilter) (string, []interface{}) {
	deleteBuilder := strings.Builder{}
	t := *new(T)
	deleteBuilder.WriteString(fmt.Sprintf("delete from %s", T.TableName(t)))
	whereString, values := buildWhere(filter, 0)
	deleteBuilder.WriteString(whereString)
	return deleteBuilder.String(), values
}

func BuildInsert[T Entity](withoutAutogenerated bool) string {
	insertQuery := strings.Builder{}
	t := *new(T)
	insertQuery.WriteString(fmt.Sprintf("insert into %s", T.TableName(t)))
	insertQuery.WriteString(fmt.Sprintf("\n\t(%s)", columnNames[T](withoutAutogenerated)))

	entityType := reflect.TypeOf(t)

	values := make([]string, 0, entityType.NumField())
	currParamIndex := 1
	for i := 0; i < entityType.NumField(); i++ {
		if entityType.Field(i).Tag.Get("not-mapped") == "true" {
			continue
		}

		if withoutAutogenerated && entityType.Field(i).Tag.Get("auto-generated") == "true" {
			continue
		}

		values = append(values, fmt.Sprintf("$%d", currParamIndex))
		currParamIndex++
	}

	insertQuery.WriteString(fmt.Sprintf("\nvalues\n\t(%s)", strings.Join(values, ", ")))

	if withoutAutogenerated {
		insertQuery.WriteString(fmt.Sprintf("\nreturning %s", columnNamesWithAliases[T](true)))
	}

	return insertQuery.String()
}

func RowsToEntities[T Entity](rows *sql.Rows) ([]T, error) {
	entities := make([]T, 0)
	entityType := reflect.TypeOf(*new(T))
	mappedFieldsIndexes := make([]int, 0)

	for i := 0; i < entityType.NumField(); i++ {
		if entityType.Field(i).Tag.Get("not-mapped") == "true" {
			continue
		} else {
			mappedFieldsIndexes = append(mappedFieldsIndexes, i)
		}
	}

	for rows.Next() {
		entity := new(T)
		entityValue := reflect.ValueOf(*entity)

		columns := make([]interface{}, 0, len(mappedFieldsIndexes))
		for _, index := range mappedFieldsIndexes {
			field := entityValue.Field(index)
			columns = append(columns, field)
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
		return fmt.Sprintf("select %s", columnNamesWithAliases[T](false))
	case Sum:
		return fmt.Sprintf("select sum(%s)", specialSelect)
	case Count:
		return fmt.Sprintf("select count(%s)", specialSelect)
	default:
		return fmt.Sprintf("select %s", specialSelect)
	}
}

func columnNames[T Entity](withoutAutogenerated bool) string {
	t := reflect.TypeOf(*new(T))
	sqlColumns := make([]string, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.Tag.Get("not-mapped") == "true" {
			continue
		}

		if withoutAutogenerated && field.Tag.Get("auto-generated") == "true" {
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

func columnNamesWithAliases[T Entity](onlyAutogenerated bool) string {
	t := reflect.TypeOf(*new(T))

	sqlColumns := make([]string, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if onlyAutogenerated && field.Tag.Get("auto-generated") != "true" {
			continue
		}

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

func buildWhere(filter any, firstParamIndex int) (string, []interface{}) {
	if filter == nil {
		return "", nil
	}
	whereBuilder := strings.Builder{}
	whereBuilder.WriteString("\n where 1 == 1")

	values := make([]interface{}, 0)
	v := reflect.ValueOf(filter)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	t := v.Type()
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
