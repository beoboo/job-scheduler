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

func (s *Stream) lock(id int) {
	//println("locking", id)
	s.m.Lock()
}

func (s *Stream) unlock(id int) {
	//println("unlocking", id)
	s.m.Unlock()
}

func (s *Stream) IsClosed() bool {
	s.lock(1)
	defer s.unlock(1)

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
	s.lock(2)
	defer s.unlock(2)

	return s.pos < len(s.lines)
}

func (s *Stream) readNext() *Line {
	s.lock(3)
	defer s.unlock(3)

	pos := s.pos
	s.pos += 1

	return &s.lines[pos]
}

func (s *Stream) Write(line Line) error {
	if s.IsClosed() {
		return io.ErrClosedPipe
	}

	s.lock(4)
	defer s.unlock(4)
	s.lines = append(s.lines, line)

	return nil
}

func (s *Stream) ResetPos() {
	s.lock(5)
	defer s.unlock(5)

	s.pos = 0
}

func (s *Stream) Close() {
	s.lock(6)
	defer s.unlock(6)

	if s.closed {
		return
	}

	close(s.close)

	s.closed = true
}
