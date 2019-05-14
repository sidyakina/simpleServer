package use_cases

import (
	"github.com/sidyakina/simpleServer/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestPairs struct {
	name string
	req *domain.SortingRequest
	wRes domain.SortingResponse
}

func TestSorter_Sort(t *testing.T) {
	tests := []TestPairs{
		{name: "unique false",
			req: &domain.SortingRequest{Array:[]int{1, 7, 1, 7, 0}, Unique:false},
			wRes:domain.SortingResponse{Array:[]int{0, 1, 1, 7, 7}}},
		{name: "unique true",
			req: &domain.SortingRequest{Array:[]int{1, 7, 1, 7, 0}, Unique:true},
			wRes:domain.SortingResponse{Array:[]int{0, 1, 7}}},
	}

	for _, pair := range tests {
		t.Run(pair.name, func(t *testing.T) {
			resp := Sorter{}.Sort(pair.req)
			assert.Equal(t, pair.wRes, resp, pair.name)
		})
	}

}