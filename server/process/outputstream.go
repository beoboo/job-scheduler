package process

import (
	"github.com/beoboo/job-worker-service/protocol"
)

type OutputStream struct {
	channel string
	time    int64
	text    string
}

func (os *OutputStream) ToProtocol() protocol.OutputStream {
	return protocol.OutputStream{
		Channel: os.channel,
		Time:    int(os.time),
		Text:    os.text,
	}
}
