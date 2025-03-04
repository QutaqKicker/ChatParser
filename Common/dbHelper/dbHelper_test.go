package dbHelper

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func queriesIsSame(firstQuery, secondQuery string) bool {
	clearQuery := func(query string) string {
		query = strings.ReplaceAll(query, "\n", "")
		query = strings.ReplaceAll(query, "\r", "")
		query = strings.ReplaceAll(query, "\t", "")
		query = strings.ReplaceAll(query, " as ", "")

		bytes := []byte(query)
		offset := 0
		for i := 0; i < len(bytes); i++ {
			if bytes[i] == ' ' {
				offset++
				continue
			}

			if offset > 0 {
				bytes[i-offset] = bytes[i]
			}
		}
		return string(bytes[:len(bytes)-offset])
	}

	firstQuery, secondQuery = clearQuery(firstQuery), clearQuery(secondQuery)
	return strings.EqualFold(firstQuery, secondQuery)
}

type TestEntity struct {
	Id              int32  `auto-generated:"true"`
	MappedColumn    int32  `column:"mapped_column"`
	NotMappedColumn string `not-mapped:"true"`
	Message         string
	Created         time.Time
}

func (t TestEntity) TableName() string {
	return "test_entities"
}

func TestBuildSelect(t *testing.T) {
	tests := []struct {
		selectType            SelectType
		specialSelect, result string
	}{
		{All, "", "select id, mapped_column as MappedColumn, message, created"},
		{Sum, "mapped_column", "select sum(mapped_column)"},
		{Count, "mapped_column", "select count(mapped_column)"},
		{Special, "\"SpecialSelect\"", "select \"SpecialSelect\""},
	}

	for _, test := range tests {
		query := buildSelect[TestEntity](test.selectType, test.specialSelect)
		if !queriesIsSame(query, test.result) {
			assert.Equal(t, query, test.result)
		}
	}
}

func TestColumnNames(t *testing.T) {
	tests := []struct {
		withoutAutogenerated bool
		result               string
	}{
		{false, "id, mapped_column, message, created"},
		{true, "mapped_column, message, created"},
	}
	for _, test := range tests {
		query := columnNames[TestEntity](test.withoutAutogenerated)
		if !queriesIsSame(query, test.result) {
			assert.Equal(t, query, test.result)
		}
	}
}

func TestColumnNamesWithAliases(t *testing.T) {
	tests := []struct {
		onlyAutogenerated bool
		result            string
	}{
		{false, "id, mapped_column as MappedColumn, message, created"},
		{true, "id"},
	}
	for _, test := range tests {
		query := columnNamesWithAliases[TestEntity](test.onlyAutogenerated)
		if !queriesIsSame(query, test.result) {
			assert.Equal(t, query, test.result)
		}
	}
}

type TestEntityFilter struct {
	Id             int32     `column:"id" relation:"="`
	Ids            []int32   `column:"id" relation:"in"`
	MinCreatedDate time.Time `column:"created" relation:">"`
	MaxCreatedDate time.Time `column:"created" relation:"<"`
	MappedColumn   int       `column:"mapped_column" relation:"="`
	SubMessage     string    `column:"message" relation:"like"`
}

func TestBuildWhere(t *testing.T) {
	emptyFilter := TestEntityFilter{}
	partialFilter := emptyFilter
	partialFilter.Ids = []int32{1, 2, 3}
	partialFilter.SubMessage = "test???"
	fullFilter := partialFilter
	fullFilter.Id = 123
	fullFilter.MinCreatedDate = time.Date(2025, 01, 01, 15, 00, 01, 0, time.Local)
	fullFilter.MaxCreatedDate = time.Date(2025, 02, 02, 16, 00, 02, 0, time.Local)
	fullFilter.MappedColumn = 321

	whereStart := "where 1 = 1 "
	tests := []struct {
		filter          any
		firstParamIndex int
		result          string
		params          []interface{}
	}{
		{emptyFilter, 10, "", []interface{}{}},
		{partialFilter, 1, "and id in ($1, $2, $3) and message like '%' || $4 || '%'",
			[]interface{}{int32(1), int32(2), int32(3), partialFilter.SubMessage}},
		{fullFilter, 4, "and id = $4 and id in ($5, $6, $7) and created > $8 and created < $9" +
			" and mapped_column = $10 and message like '%' || $11 || '%'",
			[]interface{}{fullFilter.Id, int32(1), int32(2), int32(3), fullFilter.MinCreatedDate, fullFilter.MaxCreatedDate,
				fullFilter.MappedColumn, partialFilter.SubMessage}},
	}

	for _, test := range tests {
		query, params := buildWhere(test.filter, test.firstParamIndex)
		if !queriesIsSame(query, whereStart+test.result) {
			t.Error(query, whereStart+test.result)
		}
		for i, param := range params {
			if param != test.params[i] {
				assert.Equal(t, query, test.result)
			}
		}
	}
}

func TestBuildSorter(t *testing.T) {
	emptySort := make([]SortField, 0, 3)
	singleSort := append(emptySort, SortField{FieldName: "id", Direction: Asc})
	tripleSort := append(singleSort, SortField{FieldName: "created", Direction: Asc}, SortField{FieldName: "mapped_column", Direction: Desc})

	tests := []struct {
		sortFields []SortField
		result     string
	}{
		{emptySort, "order by created desc"},
		{singleSort, "order by id asc"},
		{tripleSort, "order by id asc, created asc, mapped_column desc"},
	}

	for _, test := range tests {
		query := buildSorter(test.sortFields)
		if !queriesIsSame(query, test.result) {
			assert.Equal(t, query, test.result)
		}
	}
}

func TestBuildInsert(t *testing.T) {
	tests := []struct {
		withoutAutogenerated bool
		result               string
	}{
		{false, "insert into test_entities (id, mapped_column, message, created) values ($1, $2, $3, $4)"},
		{true, "insert into test_entities (mapped_column, message, created) values ($1, $2, $3) returning id"},
	}

	for _, test := range tests {
		query := BuildInsert[TestEntity](test.withoutAutogenerated)
		if !queriesIsSame(query, test.result) {
			assert.Equal(t, query, test.result)
		}
	}
}

func TestBuildDelete(t *testing.T) {
	filter := TestEntityFilter{Id: 123}
	requiredQuery := "delete from test_entities where 1 = 1 and id = $1"
	query, params := BuildDelete[TestEntity](filter)
	if !queriesIsSame(query, requiredQuery) {
		assert.Equal(t, query, requiredQuery)
	}
	assert.Equal(t, params[0], filter.Id)
}

func TestBuildUpdate(t *testing.T) {
	singleUpdate := SetUpdate("mapped_column", int32(111))
	tripleUpdate := singleUpdate.AndUpdate("id", int32(321)).AndUpdate("Message", "new message")
	doubleFilter := TestEntityFilter{Id: 123, SubMessage: "test"}
	singleFilter := TestEntityFilter{Id: 123}
	tests := []struct {
		updateValues UpdateValues
		filter       any
		result       string
		params       []interface{}
	}{
		{singleUpdate, doubleFilter,
			"update test_entities set mapped_column = $1 where 1=1 and id = $2 and message like '%' || $3 || '%'",
			[]any{int32(111), doubleFilter.Id, doubleFilter.SubMessage}},
		{tripleUpdate, singleFilter,
			"update test_entities set  mapped_column = $1, id = $2, message = $3 where 1=1 and id = $4",
			[]any{int32(111), int32(321), "new message", singleFilter.Id}},
	}
	for _, test := range tests {
		query, params := BuildUpdate[TestEntity](test.updateValues, test.filter)
		if !queriesIsSame(query, test.result) {
			t.Error(query, test.result)
		}

		for i, param := range params {
			assert.Equal(t, param, test.params[i])
		}
	}
}

func TestBuildQuery(t *testing.T) {
	emptyRequest := SelectBuildRequest{}

	requestWithFilter := emptyRequest
	requestWithFilter.WithFilter(TestEntityFilter{Id: 123})

	requestWithSorts := requestWithFilter
	requestWithSorts.WithSorts([]SortField{{FieldName: "id", Direction: Desc}})

	requestWithSelectType := requestWithSorts
	requestWithSelectType.SetSelectType(Special, "'SpecialSelect'")

	requestWithTakeAndSkip := requestWithSelectType
	requestWithTakeAndSkip.NeedTake(20).NeedSkip(100)

	requestWithOnlyTakeAndSkip := SelectBuildRequest{Take: 2, Skip: 4}

	tests := []struct {
		request SelectBuildRequest
		result  string
		params  []any
	}{
		{emptyRequest,
			"select id, mapped_column as mappedcolumn, message, created from test_entities order by created desc",
			[]any{}},
		{requestWithFilter,
			"select id, mapped_column as mappedcolumn, message, created from test_entities where 1 = 1 and id = $1 order by created desc",
			[]any{int32(123)}},
		{requestWithSorts, "select id, mapped_column as mappedcolumn, message, created from test_entities where 1 = 1 and id = $1 order by id desc", []any{int32(123)}},
		{requestWithSelectType, "select 'SpecialSelect' from test_entities where 1 = 1 and id = $1 order by id desc", []any{int32(123)}},
		{requestWithTakeAndSkip, "select 'SpecialSelect' from test_entities where 1 = 1 and id = $1 order by id desc limit 20 offset 100", []any{int32(123)}},
		{requestWithOnlyTakeAndSkip, "select id, mapped_column as mappedcolumn, message, created from test_entities order by created desc LIMIT 2 OFFSET 4", []any{}},
	}

	for _, test := range tests {
		query, params := BuildQuery[TestEntity](test.request)
		if !queriesIsSame(query, test.result) {
			assert.Equal(t, query, test.result)
		}

		for i, param := range params {
			assert.Equal(t, param, test.params[i])
		}
	}
}
