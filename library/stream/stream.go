package stream

import (
	"io"
	"sync"
	"time"
)

type Stream struct {
	lines  Lines
	pos    int
	close  chan bool
	closed bool
	m      sync.Mutex
}

func NewStream() *Stream {
	s := &Stream{
		lines: Lines{},
		close: make(chan bool, 1),
	}

	return s
}

func (s *Stream) IsClosed() bool {
	s.m.Lock()
	defer s.m.Unlock()

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
	s.m.Lock()
	defer s.m.Unlock()

	return s.pos < len(s.lines)
}

func (s *Stream) readNext() *Line {
	s.m.Lock()
	defer s.m.Unlock()

	pos := s.pos
	s.pos += 1

	return &s.lines[pos]
}

func (s *Stream) Write(line Line) error {
	if s.IsClosed() {
		return io.ErrClosedPipe
	}

	s.m.Lock()
	defer s.m.Unlock()
	s.lines = append(s.lines, line)

	return nil
}

func (s *Stream) ResetPos() {
	s.m.Lock()
	defer s.m.Unlock()

	s.pos = 0
}

func (s *Stream) Close() {
	s.m.Lock()
	defer s.m.Unlock()

	close(s.close)

	s.closed = true
}
