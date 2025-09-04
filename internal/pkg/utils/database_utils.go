package utils

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// function that generate query without explicitly writing a field name
// lets keep it flexible so interface can be student , exec ..any
func GenerateInsertQuery(tableName string, model interface{}) string {
	modelType := reflect.TypeOf(model)
	var columns string
	placeholders := []string{}
	paramCount := 1
	for i := 0; i < modelType.NumField(); i++ {
		//access the db tag
		dbTag := modelType.Field(i).Tag.Get("db")
		fmt.Println("dbTag:", dbTag)
		//extract the column
		dbTag = strings.TrimSuffix(dbTag, ",omitempty")
		//condition we dont need the id from fields we need the rest of the columns we are not posting the id posting rest of the columns , id is autmatically generated
		if dbTag != "" && dbTag != "id" && dbTag != "created_at" && dbTag != "updated_at" { //skip the ID field if its auto increment
			if columns != "" {
				columns += ", "
				//placeholders += ", " ,,pgx use $ not ?  like sql

			}
			columns += dbTag
			placeholders = append(placeholders, fmt.Sprintf("$%d", paramCount))
			paramCount++
		}
	}
	fmt.Printf("INSERT INTO %s (%s) VALUES (%s)\n", tableName, columns, strings.Join(placeholders, ", "))
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, columns, strings.Join(placeholders, ", "))
}

// return them in a list we dont know how many values
func GetStructValues(model interface{}) []interface{} {

	modelValue := reflect.ValueOf(model)
	modelType := modelValue.Type()

	values := []interface{}{}
	for i := 0; i < modelType.NumField(); i++ {

		dbTag := modelType.Field(i).Tag.Get("db")
		if dbTag != "" && dbTag != "id,omitempty" && dbTag != "created_at,omitempty" && dbTag != "updated_at,omitempty" {
			values = append(values, modelValue.Field(i).Interface())

		}

	}
	log.Println("Values:", values)
	return values
}
func GetPaginationParams(r *http.Request) (int, int) {

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		//fmt.Println("default ")

		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		//fmt.Println("errrrrr")
		limit = 10
	}
	return page, limit

}
func isValidSortField(field string) bool {
	validFields := map[string]bool{
		"word":     true,
		"phonetic": true,
		"origin":   true,
	}
	return validFields[field]
}

//func AddFilters(r *http.Request, query string, args []interface{}) (string, []interface{}) {
//	params := map[string]string{
//		"word":     "word",
//		"phonetic": "phonetic",
//		"origin":   "origin",
//	}
//	for param, dbField := range params {
//		value := r.URL.Query().Get(param)
//		if value != "" {
//			query += " AND " + dbField + " ILIKE $%d" // ILIKE for case-insensitive
//			args = append(args, "%"+value+"%")
//		}
//	}
//	// General search on word....../words?search=runjhj
//	search := r.URL.Query().Get("search")
//	if search != "" {
//		query += " AND word ILIKE $%d"
//		args = append(args, "%"+search+"%")
//	}
//	fmt.Println(len(args))
//	return fmt.Sprintf(query, len(args)+1), args
//}

// ///words?sortBy=word:asc&sortBy=origin:desc
func AddSorting(r *http.Request, query string) string {
	sortParams := r.URL.Query()["sortBy"]
	if len(sortParams) > 0 {
		query += " ORDER BY"
		for i, params := range sortParams {
			parts := strings.Split(params, ":")
			if len(parts) != 2 {
				continue
			}
			field, order := parts[0], parts[1]
			if !isValidSortField(field) || !isValidSortOrder(order) {
				continue
			}
			if i > 0 {
				query += ","
			}
			//field = "word" , order = "asc"

			query += " " + field + " " + order
		}
	}
	return query
}

func isValidSortOrder(order string) bool {
	return order == "asc" || order == "desc"
}

func AddFilters(r *http.Request, query string, args []interface{}) (string, []interface{}) {
	paramCount := len(args) + 1 // Start from next $n
	params := map[string]string{
		"word":     "word",
		"phonetic": "phonetic",
		"origin":   "origin",
	}
	for param, dbField := range params {
		value := r.URL.Query().Get(param)
		if value != "" {
			query += fmt.Sprintf(" AND %s ILIKE $%d", dbField, paramCount)
			args = append(args, "%"+value+"%")
			paramCount++
		}
	}
	// General search on word
	search := r.URL.Query().Get("search")
	if search != "" {
		query += fmt.Sprintf(" AND word ILIKE $%d", paramCount)
		args = append(args, "%"+search+"%")
	}
	return query, args
}
