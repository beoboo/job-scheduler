package library

type StatusRequestData struct {
	Id string `json:"id"`
}

type StatusResponseData struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}
