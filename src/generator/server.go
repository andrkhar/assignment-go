package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type (
	SensorData struct {
		SensorValue string `json:"Sensor Value"`
		ID1         string `json:"ID1"`
		ID2         string `json:"ID2"`
		Timestamp   string `json:"Timestamp"`
	}
)

const (
	timeLayoutMessage  = "Mon 01/02/2006-15:04:05"
	messageLocation    = "Asia/Singapore"
	timeLayoutDatabase = "2006-01-02 03:04:05"
)

var locationSG, _ = time.LoadLocation(messageLocation)
var locationUTC, _ = time.LoadLocation("UTC")

var timer *time.Timer

func getRandomData() (data *SensorData) {
	data = new(SensorData)
	data.SensorValue = strconv.Itoa(rand.Intn(101))
	data.ID1 = strconv.Itoa(rand.Intn(5) + 1)
	data.ID2 = string(int('A') + rand.Intn(6))
	data.Timestamp = time.Now().Format(timeLayoutMessage)
	return
}

func sendData() {
	data := getRandomData()
	fmt.Println("Data sent: ", data)
	var duration time.Duration = 2 * time.Second
	timer.Reset(duration)
}

func start(c echo.Context) error {
	timer.Stop()
	timer.Reset(0)
	return c.String(http.StatusOK, "Generation started.")
}

func stop(c echo.Context) error {
	timer.Stop()
	return c.String(http.StatusOK, "Generation stopped.")
}

func main() {
	rand.Seed(0)
	timer = time.AfterFunc(time.Hour, sendData)
	timer.Stop()
	e := echo.New()
	e.POST("/start", start)
	e.POST("/stop", stop)
	e.Logger.Fatal(e.Start(":1324"))
}
