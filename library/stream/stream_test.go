package stream

import (
	"io"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	s := NewStream()
	defer s.Close()

	if s.IsClosed() {
		t.Fatalf("Stream should be open")
	}
}

func write(s *Stream, t string) {
	go func() {
		_ = s.Write(buildLine(t))
	}()
}

func TestWrite(t *testing.T) {
	s := NewStream()
	defer s.Close()

	expected := "line"
	write(s, expected)

	assertRead(t, s, expected)
}

func TestResetPos(t *testing.T) {
	s := NewStream()
	defer s.Close()

	expected := "line"
	write(s, expected)

	assertRead(t, s, expected)

	s.ResetPos()

	assertRead(t, s, expected)
}

func TestAllLines(t *testing.T) {
	s := NewStream()

	write(s, "line1")
	write(s, "line2")

	assertRead(t, s, "line1")
	assertRead(t, s, "line2")
}

func TestWaitNext(t *testing.T) {
	s := NewStream()

	go func() {
		time.Sleep(100 * time.Millisecond)
		s.Close()
	}()

	assertError(t, s, io.EOF)
}

func TestCannotWriteToClosedStream(t *testing.T) {
	s := NewStream()
	s.Close()

	errCh := make(chan error, 1)
	go func() {
		time.Sleep(10 * time.Millisecond)
		err := s.Write(buildLine("invalid"))

		errCh <- err
	}()

	if <-errCh == nil {
		t.Fatal("Should not be able to write if stream is closed")
	}
}

func TestClose(t *testing.T) {
	s := NewStream()
	s.Close()

	assertError(t, s, io.EOF)
	assertError(t, s, io.EOF)
}

func buildLine(t string) Line {
	return Line{
		Text: t,
		Time: time.Now(),
	}
}

func assertRead(t *testing.T, s *Stream, expected string) {
	l, err := s.Read()

	if l.Text != expected {
		t.Fatalf("Didn't read successfully, expected \"%s\", got \"%s\"", expected, l.Text)
	}

	if err != nil {
		t.Fatalf("Didn't expect any error, got \"%v\"", err)
	}
}

func assertError(t *testing.T, s *Stream, expected error) {
	l, err := s.Read()

	if l != nil {
		t.Fatalf("Didn't expect anything read, got \"%s\"", l)
	}

	if err == nil {
		t.Fatalf("Expected error %v", expected)
	}
}