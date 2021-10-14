package viewmodel

type ResponseVm struct {
	Body       RespBodyVm
	StatusCode int
}

type RespBodyVm struct {
	Message    interface{} `json:"message"`
	DataVm     interface{} `json:"data"`
	Pagination interface{} `json:"pagination"`
}
