package protocol

type StopRequestData struct {
	Pid int `json:"pid"`
}

type StopResponseData struct {
	Pid    int    `json:"pid"`
	Status string `json:"status"`
}
