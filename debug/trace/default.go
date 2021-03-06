package trace

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/micro/go-micro/util/ring"
)

type trace struct {
	opts Options

	// ring buffer of traces
	buffer *ring.Buffer
}

func (t *trace) Read(opts ...ReadOption) ([]*Span, error) {
	return []*Span{}, nil
}

func (t *trace) Start(ctx context.Context, name string) *Span {
	span := &Span{
		Name:     name,
		Trace:    uuid.New().String(),
		Id:       uuid.New().String(),
		Started:  time.Now(),
		Metadata: make(map[string]string),
	}

	// return span if no context
	if ctx == nil {
		return span
	}

	s, ok := FromContext(ctx)
	if !ok {
		return span
	}

	// set trace id
	span.Trace = s.Trace
	// set parent
	span.Parent = s.Id

	// return the sapn
	return span
}

func (t *trace) Finish(s *Span) error {
	// set finished time
	s.Finished = time.Now()

	// save the span
	t.buffer.Put(s)

	return nil
}

func NewTrace(opts ...Option) Trace {
	var options Options
	for _, o := range opts {
		o(&options)
	}

	return &trace{
		opts: options,
		// the last 64 requests
		buffer: ring.New(64),
	}
}
