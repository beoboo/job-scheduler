package process

import (
	"github.com/beoboo/job-worker-service/protocol"
)

type OutputStream struct {
	Channel string
	Time    int64
	Text    string
}

func (os *OutputStream) ToProtocol() protocol.OutputStream {
	return protocol.OutputStream{
		Channel: os.Channel,
		Time:    int(os.Time),
		Text:    os.Text,
	}
}
