package library

type StartRequestData struct {
	Executable string `json:"executable"`
	Args       string `json:"args"`
}

type StartResponseData struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}
