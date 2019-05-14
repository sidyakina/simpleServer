package main

import (
	"github.com/sidyakina/simpleServer/adapters"
	"github.com/sidyakina/simpleServer/use_cases"
)

func main(){
	weather := use_cases.InitWeather()
	h := adapters.InitHandlers(use_cases.GetTime{}, use_cases.Sorter{}, weather)
	s := adapters.InitServer(h)
	s.Start("3333")
}
