package main

import (
	"fmt"
	"reflect"
	"strings"
)

type Row struct {
	ColumnA string `f:"column_a"`
	ColumnB int `f:"column_b"`
	ColumnC bool `f:"column_c"`
}

func main() {
	rows := []Row{
		{
			ColumnA: "A",
			ColumnB: 1,
			ColumnC: false,
		},
		{
			ColumnA: "B",
			ColumnB: 2,
			ColumnC: true,
		},
	}

	query := "INSERT INTO row (columnA, columnB, columnC) VALUES ($1, $2, $3) , ($4, $5, $6)"
	fmt.Println(query)

	fmt.Println(insertMany(rows))
}

func insertMany(rows interface{}) string {
	queryParts := []string{"INSERT INTO"}

	baseType := reflect.TypeOf(rows)
	table := baseType.Elem()

	queryParts = append(queryParts, strings.ToLower(table.Name()))

	var tableFields []string
	for i:= 0; i < table.NumField()-1; i++ {
		tableFields = append(tableFields, table.Field(i).Tag.Get("f")+",")
	}
	tableFields = append(tableFields, table.Field(table.NumField()-1).Tag.Get("f"))
	tableFields = append([]string{"("}, tableFields...)
	tableFields = append(tableFields, ")")

	queryParts = append(queryParts, strings.Join(tableFields, " "))
	queryParts = append(queryParts, "VALUES")

	baseValue := reflect.ValueOf(rows)

	var values []string
	for i:=0; i< baseValue.Len(); i++ {
		var partialValues []string
		row := baseValue.Index(i)
		for j:=0; j<row.NumField()-1; j++ {
			partialValues = append(partialValues, fmt.Sprintf("$%d,",((i*row.NumField())+j)+1))
		}
		partialValues = append(partialValues, fmt.Sprintf("$%d",((i*row.NumField())+row.NumField()-1)+1))
		partialValues = append([]string{"("}, partialValues...)
		partialValues = append(partialValues, ")")
		values = append(values, strings.Join(partialValues, " "))
	}

	queryParts = append(queryParts, strings.Join(values, ", "))
	return strings.Join(queryParts, " ")
}
