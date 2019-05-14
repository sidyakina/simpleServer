package domain

type SortingRequest struct {
	Array []int `json:"array"`
	Unique bool `json:"uniq"`
}

type SortingResponse struct {
	Array []int `json:"array"`
}
