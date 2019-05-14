package main

import (
	"github.com/sidyakina/simpleServer/adapters"
	"github.com/sidyakina/simpleServer/use_cases"
)

func main(){
	h := adapters.InitHandlers(use_cases.GetTime{}, use_cases.Sorter{}, use_cases.Weather{})
	s := adapters.InitServer(h)
	s.Start("3333")
}
