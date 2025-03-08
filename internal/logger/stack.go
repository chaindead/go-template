package logger

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type state struct {
	b []byte
}

// Write implement fmt.Formatter interface.
func (s *state) Write(b []byte) (n int, err error) {
	s.b = b

	return len(b), nil
}

// Width implement fmt.Formatter interface.
func (s *state) Width() (wid int, ok bool) {
	return 0, false
}

// Precision implement fmt.Formatter interface.
func (s *state) Precision() (prec int, ok bool) {
	return 0, false
}

// Flag implement fmt.Formatter interface.
func (s *state) Flag(_ int) bool {
	return false
}

func frameField(f errors.Frame, s *state, c rune) string {
	f.Format(s, c)

	return string(s.b)
}

// marshalStack modified version of github.com/rs/zerolog@v1.33.0/pkgerrors/stacktrace.go
// for human-friendly/shorter stack marshaling (json in original)
//
// implements pkg/errors stack trace marshaling.
//
// zerolog.ErrorStackMarshaler = marshalStack
func marshalStack(err error) interface{} {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}
	var sterr stackTracer
	var ok bool
	for err != nil {
		sterr, ok = err.(stackTracer)
		if ok {
			break
		}

		u, ok := err.(interface {
			Unwrap() error
		})
		if !ok {
			return nil
		}

		err = u.Unwrap()
	}
	if sterr == nil {
		return nil
	}

	st := sterr.StackTrace()
	if len(st) > 2 {
		st = st[:len(st)-2]
	}

	s := &state{}
	out := make([]string, 0, len(st))
	for _, frame := range st {
		out = append(out,
			fmt.Sprintf("%s:%s(%s)",
				frameField(frame, s, 's'),
				frameField(frame, s, 'd'),
				frameField(frame, s, 'n'),
			),
		)
	}

	return strings.Join(out, "; ")
}
