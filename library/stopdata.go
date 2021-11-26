package library

type StopRequestData struct {
	Id string `json:"id"`
}

type StopResponseData struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}
