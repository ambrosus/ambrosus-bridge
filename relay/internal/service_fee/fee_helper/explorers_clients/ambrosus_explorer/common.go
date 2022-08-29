package ambrosus_explorer

type Meta struct {
	Code    int    `json:"code"`
	Message string `json:"message"` // when request is unsuccessful
}

type Pagination struct {
	HasNext     bool `json:"hasNext"`
	HasPrevious bool `json:"hasPrevious"`
	Previous    int  `json:"previous"`
}
