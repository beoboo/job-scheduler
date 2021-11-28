package stream

import (
	"github.com/beoboo/job-scheduler/library/log"
	"io"
	"sync"
	"time"
)

type Stream struct {
	logger *log.Logger
	lines  Lines
	pos    int
	close  chan bool
	closed bool
	m      sync.Mutex
}

func New(logger *log.Logger) *Stream {
	s := &Stream{
		logger: logger,
		lines:  Lines{},
		close:  make(chan bool, 1),
	}

	return s
}

func (s *Stream) IsClosed() bool {
	s.lock("IsClosed")
	defer s.unlock("IsClosed")

	return s.closed
}

func (s *Stream) Read() (*Line, error) {
	for {
		if s.hasData() {
			return s.readNext(), nil
		} else {
			select {
			case <-time.After(100 * time.Millisecond):
				continue
			case <-s.close:
			}
			break
		}
	}

	return nil, io.EOF
}

func (s *Stream) hasData() bool {
	s.lock("hasData")
	defer s.unlock("hasData")

	return s.pos < len(s.lines)
}

func (s *Stream) readNext() *Line {
	s.lock("readNext")
	defer s.unlock("readNext")

	pos := s.pos
	s.pos += 1

	return &s.lines[pos]
}

func (s *Stream) Write(line Line) error {
	if s.IsClosed() {
		return io.ErrClosedPipe
	}

	s.lock("Write")
	defer s.unlock("Write")
	s.lines = append(s.lines, line)

	return nil
}

func (s *Stream) ResetPos() {
	s.lock("ResetPos")
	defer s.unlock("ResetPos")

	s.pos = 0
}

func (s *Stream) Close() {
	s.lock("Close")
	defer s.unlock("Close")

	if s.closed {
		return
	}

	close(s.close)

	s.closed = true
}

func (s *Stream) lock(id string) {
	s.debug("Stream locking %s", id)
	s.m.Lock()
}

func (s *Stream) unlock(id string) {
	s.debug("Stream unlocking %s", id)
	s.m.Unlock()
}

func (s *Stream) debug(format string, args ...interface{}) {
	if s.logger != nil {
		s.logger.Debugf(format+"\n", args...)
	}
}
