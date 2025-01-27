package tests

import (
	"testing"
	"time"
)

type testEntity struct {
	Id              int    `auto-generated:"true"`
	MappedColumn    int    `column:"mapped_column"`
	NotMappedColumn string `not-mapped:"true"`
	Message         string
	Created         time.Time
}

func (t *testEntity) TableName() string {
	return "test_entities"
}

type testEntityFilter struct {
	Id             int       `column:"id" relation:"="`
	Ids            []int     `column:"id" relation:"in"`
	MinCreatedDate time.Time `column:"created" relation:"<"`
	MaxCreatedDate time.Time `column:"created" relation:">"`
	MappedColumn   int       `column:"mapped_column" relation:"="`
	SubMessage     string    `column:"message" relation:"like"`
}

func TestBuildSelect(t *testing.T) {

}

func TestColumnNames(t *testing.T) {

}

func TestColumnNamesWithAliases(t *testing.T) {

}

func TestBuildWhere(t *testing.T) {

}

func TestBuildSorter(t *testing.T) {

}

func TestBuildInsert(t *testing.T) {

}

func TestBuildDelete(t *testing.T) {

}

func TestBuildUpdate(t *testing.T) {

}

func TestBuildQuery(t *testing.T) {

}
