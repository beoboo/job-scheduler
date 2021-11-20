package protocol

type StartData struct {
	Command string `json:"command"`
	Args    string `json:"args"`
}
