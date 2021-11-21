package protocol

type OutputRequestData struct {
	Id string `json:"id"`
}

type OutputStream struct {
	Channel string `json:"channel"`

	// TODO: This should consider the server timezone
	Time int    `json:"time"`
	Text string `json:"text"`
}

type OutputResponseData struct {
	Id     string         `json:"id"`
	Output []OutputStream `json:"output"`
}
