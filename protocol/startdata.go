package protocol

type StartRequestData struct {
	Executable string `json:"executable"`
	Args       string `json:"args"`
}

type StartResponseData struct {
	Pid    int    `json:"pid"`
	Status string `json:"status"`
}
