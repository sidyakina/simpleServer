package use_cases

import (
	"github.com/sidyakina/simpleServer/domain"
	"sort"
)


type Sorter struct {}

func (s Sorter) Sort (request *domain.SortingRequest) domain.SortingResponse {
	sort.Ints(request.Array)
	resp := domain.SortingResponse{Array: request.Array}
	if request.Unique {
		newArray := make([]int, 0, len(request.Array))
		newArray = append(newArray, request.Array[0])
		lastAddedElement := request.Array[0]
		for i := 1; i < len(request.Array); i++ {
			if request.Array[i] != lastAddedElement {
				newArray = append(newArray, request.Array[i])
				lastAddedElement = request.Array[i]
			}
		}
		resp.Array = newArray
	}
	return resp
}
