package db_handler

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func CreateInsertQuery(q interface{}) (*string, error) {
	if reflect.ValueOf(q).Kind() == reflect.Struct {
		t := reflect.TypeOf(q).Name()
		v := reflect.ValueOf(q)

		var queryNames string
		var queryValues string
		firstElem := true
		for i := 0; i < v.NumField(); i++ {
			if v.Type().Field(i).Name == "ID" {
				continue //skipping due to autoincrement field
			}
			switch v.Field(i).Kind() {
			case reflect.Int:
				if firstElem {
					queryValues = fmt.Sprintf("%s%d", queryValues, v.Field(i).Int())
					queryNames = fmt.Sprintf("%s%d", queryNames, v.Type().Field(i).Name)
				} else {
					queryValues = fmt.Sprintf("%s, %d", queryValues, v.Field(i).Int())
					queryNames = fmt.Sprintf("%s, %d", queryNames, v.Type().Field(i).Name)
				}
			case reflect.String:
				if firstElem {
					queryValues = fmt.Sprintf("%s\"%s\"", queryValues, v.Field(i).String())
					queryNames = fmt.Sprintf("%s\"%s\"", queryNames, strings.ToLower(v.Type().Field(i).Name))
				} else {
					queryValues = fmt.Sprintf("%s, \"%s\"", queryValues, v.Field(i).String())
					queryNames = fmt.Sprintf("%s, \"%s\"", queryNames, strings.ToLower(v.Type().Field(i).Name))
				}
			default:
				return nil, errors.New("unsupported type")
			}
			firstElem = false
		}
		query := fmt.Sprintf("insert into %s (%s) values (%s)", strings.ToLower(t), queryNames, queryValues)
		fmt.Println(query)
		return &query, nil
	}
	return nil, errors.New("unsupported type")
}
