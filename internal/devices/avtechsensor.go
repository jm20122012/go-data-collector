package devices

import (
	"encoding/json"
	"fmt"
	"go-data-collector/internal/db"
	"go-data-collector/internal/utils"
	"net/http"
	"strconv"
)

type AvtechResponse struct {
	Sensor []AvtechData `json:"sensor"`
}

type AvtechData struct {
	Label string `json:"label"`
	TempF string `json:"tempf"`
	TempC string `json:"tempc"`
	HighF string `json:"highf"`
	HighC string `json:"highc"`
	LowF  string `json:"lowf"`
	LowC  string `json:"lowc"`
}

type AvtechSensor struct {
	DeviceConfig
}

func (a *AvtechSensor) FetchData() {
	a.Logger.Debug("fetching avtech sensor data", "deviceName", a.Name, "deviceIP", a.IP)
	var response AvtechResponse

	url := fmt.Sprintf("http://%s/getData.json", a.IP)
	resp, err := http.Get(url)
	if err != nil {
		a.Logger.Error("error performing http request for avtech sensor", "deviceName", a.Name, "deviceIP", a.IP, "error", err)
		return
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		a.Logger.Error("error decoding json for avtech sensor", "deviceName", a.Name, "deviceIP", a.IP, "error", err)
		return
	}
	a.processAvtechResponse(response)

}

func (a *AvtechSensor) processAvtechResponse(response AvtechResponse) {
	a.Logger.Debug("Processing avtech sensor response")

	var tempF float64
	var tempC float64
	tempF, err := strconv.ParseFloat(response.Sensor[0].TempF, 32)
	if err != nil {
		a.Logger.Error("error converting avtech temp f value to float", "error", err)
	}
	tempC, err = strconv.ParseFloat(response.Sensor[0].TempC, 32)
	if err != nil {
		a.Logger.Error("error converting avtech temp c value to float", "error", err)
	}

	timestamp := utils.GetUTCTimestamp()
	params := db.WriteAvtechRecordParams{
		Timestamp:    timestamp,
		TempF:        roundFloat(tempF, 2),
		TempC:        roundFloat(tempC, 2),
		DeviceID:     int32(a.DeviceID),
		DeviceTypeID: int32(a.DeviceTypeID),
	}

	a.Logger.Debug("avtech write params built", "params", params)

	// err = a.DBStore.WriteAvtechRecord(a.Ctx, params)
	// if err != nil {
	// 	a.Logger.Error("error writing avtech record", "error", err)
	// }

}
