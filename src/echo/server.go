package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	//"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

type (
	mysqlCredentials struct {
		database string
		user     string
		password string
	}
)

var sqlCred mysqlCredentials

type (
	SensorData struct {
		SensorValue string `json:"Sensor Value"`
		ID1         string `json:"ID1"`
		ID2         string `json:"ID2"`
		Timestamp   string `json:"Timestamp"`
	}
)

type (
	QueryData struct {
		ID1            string `query:"ID1"`
		ID2            string `query:"ID2"`
		StartTimestamp string `query:"start_timestamp"`
		EndTimestamp   string `query:"end_timestamp"`
	}
)

const (
	dataSourceTemplate = "user:password@tcp(mysql)/database"
	sqlQueryInsert     = "INSERT INTO data (`Sensor Value`, ID1, ID2, Timestamp) VALUES (?, ?, ?, ?)"
	timeLayoutMessage  = "Mon 01/02/2006-03:04:05"
	messageLocation    = "Asia/Singapore"
	timeLayoutDatabase = "2006-01-02 03:04:05"
	sqlQuerySelect     = "SELECT `Sensor Value`, ID1, ID2, Timestamp FROM data"
)

var locationSG, _ = time.LoadLocation(messageLocation)
var locationUTC, _ = time.LoadLocation("UTC")

func getDataSource() string {
	return fmt.Sprint(
		sqlCred.user,
		":",
		sqlCred.password,
		"@tcp(mysql)/",
		sqlCred.database,
	)
}

func store(c echo.Context) error {
	data := new(SensorData)
	if err := c.Bind(data); err != nil {
		return c.String(http.StatusBadRequest, "Error 400. Failed parse input.")
	}

	parsedTime, err := time.ParseInLocation(timeLayoutMessage, data.Timestamp, locationSG)
	if err != nil {
		return c.String(http.StatusBadRequest, "Error 400. Failed parse timestamp.")
	}

	// fmt.Println(parsedTime.Unix())
	// fmt.Println(parsedTime.Location())

	formatedTime := parsedTime.Format(timeLayoutDatabase)

	db, err := sql.Open("mysql", getDataSource())
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error 500. Failed open Database.")
	}
	defer db.Close()

	stmt, err := db.Prepare(sqlQueryInsert)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error 500. Failed prepare data for Database.")
	}
	defer stmt.Close()

	if _, err := stmt.Exec(
		data.SensorValue,
		data.ID1,
		data.ID2,
		formatedTime,
	); err != nil {
		return c.String(http.StatusInternalServerError, "Error 500. Failed to write data to Database.")
	}

	return c.JSON(http.StatusOK, data)
}

func getCustomSQLQuerySelect(data *QueryData, c echo.Context) (string, error) {
	query := sqlQuerySelect
	var filters []string

	if data.ID1 != "" {
		filter := "ID1='" + data.ID1 + "'"
		filters = append(filters, filter)
	}

	if data.ID2 != "" {
		filter := "ID2='" + data.ID2 + "'"
		filters = append(filters, filter)
	}

	if data.StartTimestamp != "" {
		i, err := strconv.ParseInt(data.StartTimestamp, 10, 64)
		if err != nil {
			return "", c.String(http.StatusBadRequest, "Error 400. Failed to parse StartTimestamp.")
		}
		t := time.Unix(i, 0)

		// fmt.Println(t.Unix())
		// fmt.Println(t.Location())

		data.StartTimestamp = t.Format(timeLayoutDatabase)
		filter := "Timestamp>='" + data.StartTimestamp + "'"

		filters = append(filters, filter)
	}

	if data.EndTimestamp != "" {
		i, err := strconv.ParseInt(data.EndTimestamp, 10, 64)
		if err != nil {
			return "", c.String(http.StatusBadRequest, "Error 400. Failed to parse EndTimestamp.")
		}
		t := time.Unix(i, 0)

		// fmt.Println(t.Unix())
		// fmt.Println(t.Location())

		data.EndTimestamp = t.Format(timeLayoutDatabase)
		filter := "Timestamp<='" + data.EndTimestamp + "'"

		filters = append(filters, filter)
	}

	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	// fmt.Println(query)
	return query, nil
}

func retrieve(c echo.Context) error {
	data := new(QueryData)
	if err := c.Bind(data); err != nil {
		return c.String(http.StatusBadRequest, "Error 400. Failed parce input.")
	}
	db, err := sql.Open("mysql", getDataSource())
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error 500. Failed open Database.")
	}
	defer db.Close()

	sqlQuery, err := getCustomSQLQuerySelect(data, c)
	if err != nil {
		return err
	}

	var records []SensorData
	rows, err := db.Query(sqlQuery)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error 500. Failed query Database.")
	}
	defer rows.Close()

	for rows.Next() {
		var record SensorData
		if err := rows.Scan(&record.SensorValue, &record.ID1, &record.ID2, &record.Timestamp); err != nil {
			return c.String(http.StatusInternalServerError, "Error 500. Failed scan a Database record.")
		}
		t, err := time.Parse(timeLayoutDatabase, record.Timestamp)
		if err != nil {
			return c.String(http.StatusBadRequest, "Error 500. Failed parse a Database record Timestamp.")
		}

		record.Timestamp = t.Format(timeLayoutMessage)
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return c.String(http.StatusInternalServerError, "Error 500. Rows error.")
	}

	return c.JSON(http.StatusOK, records)
}

func setSQLCred() {
	sqlCred.database = os.Getenv("MYSQL_DATABASE")
	sqlCred.user = os.Getenv("MYSQL_USER")
	sqlCred.password = os.Getenv("MYSQL_PASSWORD")
}

func main() {
	setSQLCred()
	e := echo.New()
	e.POST("/data", store)
	e.GET("/data", retrieve)
	e.Logger.Fatal(e.Start(":1323"))
}
