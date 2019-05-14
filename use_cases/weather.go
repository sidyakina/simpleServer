package use_cases

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

type Weather struct {
	Apikey string
}

func InitWeather() Weather {
	return Weather{Apikey: os.Getenv("APIKEY")}
}

type WeatherParams struct {
	Main MainParams `json:"main"`
}
type MainParams struct {
	Temp float64 `json:"temp"`
}


func (w Weather)GetWeather(city string) (float64, error) {
	resp, err := http.Get("http://openweathermap.org/data/2.5/weather?q=" + city + "&appid=" + w.Apikey)
	if err != nil {
		return 0.0, err
	}
	if resp.StatusCode == 404 {
		return 0.0, errors.New("not found city")
	}
	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	var r WeatherParams
	err = decoder.Decode(&r)
	if err != nil {
		return 0.0, errors.New("not found city")
	}
	return r.Main.Temp, nil
}