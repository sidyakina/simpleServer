package adapters

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sidyakina/simpleServer/domain"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type getNowTimeMock struct {
	waitTime time.Time
}
func (g getNowTimeMock) GetNowTime() time.Time {
	return g.waitTime
}
type sortArrayMock struct {}
func (s sortArrayMock) Sort (request *domain.SortingRequest) domain.SortingResponse{
	return domain.SortingResponse{Array:[]int{1, 2, 3, 4, 5}}
}
type getWeatherMock struct {
	temp float64
	err error
}
func (g getWeatherMock) GetWeather(city string) (float64, error){
	return g.temp, g.err
}

func TestSimpleServer_Start_now(t *testing.T) {
	nowMock := getNowTimeMock{waitTime: time.Now()}
	h := InitHandlers(nowMock, sortArrayMock{}, getWeatherMock{})
	s := InitServer(h)
	mockServer := httptest.NewServer(http.HandlerFunc(s.nowTime))
	defer mockServer.Close()
	res, _ := http.Get(mockServer.URL + "/now")
	response, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	assert.Equal(t, nowMock.waitTime.String(), strings.Trim(string(response), "\n"), "/now")
}

func TestSimpleServer_Start_sort(t *testing.T) {
	h := InitHandlers(getNowTimeMock{}, sortArrayMock{}, getWeatherMock{})
	s := InitServer(h)
	mockServer := httptest.NewServer(http.HandlerFunc(s.sort))
	defer mockServer.Close()
	req := domain.SortingRequest{Array:[]int{1, 2, 3, 4, 5}, Unique:true}
	bytesJSON, err := json.Marshal(req)
	res, _ := http.Post(mockServer.URL + "/sort", "application/json", bytes.NewBuffer(bytesJSON))
	var jres domain.SortingResponse
	response, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	err = json.Unmarshal(response, &jres)
	fmt.Println(err)
	assert.Equal(t, domain.SortingResponse{Array:[]int{1, 2, 3, 4, 5}}, jres, "/sort")
}

func TestSimpleServer_Start_weather(t *testing.T) {
	h := InitHandlers(getNowTimeMock{}, sortArrayMock{}, getWeatherMock{})
	s := InitServer(h)
	mockServer := httptest.NewServer(http.HandlerFunc(s.weather))
	defer mockServer.Close()
	res, _ := http.Get(mockServer.URL + "/weather")
	defer res.Body.Close()
	assert.Equal(t, 400, res.StatusCode, "/weather")
}

type testpair struct {
	err error
	temp float64
	waitCode int
	waitResp string
}

func TestSimpleServer_Start_weather2(t *testing.T) {
	tests := []testpair{
		{errors.New("not found city"), 10, 404, "not found city"},
		{nil, 10.5, 200, "10.5"},
	}
	for _, pair := range tests {
		newMock := getWeatherMock{err: pair.err, temp: pair.temp}
		h := InitHandlers(getNowTimeMock{}, sortArrayMock{}, newMock)
		s := InitServer(h)
		mockServer := httptest.NewServer(http.HandlerFunc(s.weather))
		defer mockServer.Close()
		res, _ := http.Get(mockServer.URL + "/weather?city=test")
		defer res.Body.Close()
		assert.Equal(t, pair.waitCode, res.StatusCode, "/weather")
		response, _ := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		assert.Equal(t, pair.waitResp, strings.Trim(string(response), "\n"), "/weather")

	}

}

