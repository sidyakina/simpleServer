package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sidyakina/simpleServer/domain"
	"log"
	"net/http"
	"time"
)

type SimpleServer struct {
	handlers Handlers
}

type Handlers struct {
	nowTime getNowTime
	sorting sortArray
	weather getWeather
}

type getNowTime interface {GetNowTime() time.Time}
type sortArray interface {Sort (request *domain.SortingRequest) domain.SortingResponse}
type getWeather interface {GetWeather(city string) (float64, error)}

func InitHandlers(nowTime getNowTime, sorting sortArray, weather getWeather) Handlers {
	return Handlers{nowTime: nowTime, sorting: sorting, weather: weather}
}

func InitServer(handlers Handlers) *SimpleServer {
	return &SimpleServer{handlers: handlers}
}

func (s *SimpleServer) Start (port string) {
	http.HandleFunc("/now", s.nowTime)
	http.HandleFunc("/sort", s.sort)
	http.HandleFunc("/weather", s.weather)
	err := http.ListenAndServe(":" + port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *SimpleServer) nowTime(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.NotFound(res, req)
		return
	}
	_, err := fmt.Fprintln(res, s.handlers.nowTime.GetNowTime())
	if err != nil {
		log.Println(err)
	}
}

func (s *SimpleServer) sort (res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.NotFound(res, req)
		return
	}
	var r domain.SortingRequest
	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()
	err := decoder.Decode(&r)
	if err != nil {
		log.Println(err)
		http.Error(res, "can't unmarshal request", 500)
		return
	}
	if l := len(r.Array); l == 0 || l > 100 {
		http.Error(res, "empty or too long array", 400)
		return
	}
	resp := s.handlers.sorting.Sort(&r)
	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(&resp)
	if err != nil {
		log.Println(err)
		http.Error(res, "internal error", 500)
		return
	}
	_, err = fmt.Fprintln(res, buf)
	if err != nil {
		log.Println(err)
	}
}

func (s *SimpleServer) weather (res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.NotFound(res, req)
		return
	}
	query := req.URL.Query()
	city := query.Get("city")
	if city == "" {
		http.Error(res, "empty city", 400)
		return
	}
	temp, err := s.handlers.weather.GetWeather(city)
	if err != nil && err.Error() == "not found city" {
		log.Println(err)
		http.Error(res, "not found city", 404)
		return
	}
	if err != nil {
		log.Println(err)
		http.Error(res, "internal error", 500)
		return
	}
	_, err = fmt.Fprintln(res, temp)
	if err != nil {
		log.Println(err)
	}
}