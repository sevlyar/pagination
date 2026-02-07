package pagination

import (
	"fmt"
)

// Span presents position and size of collection items chunk in the collection.
type Span struct {
	Pos  *Cursor
	size int
}

// NewFirstSpan returns a new span of the defined size and in the starting position.
func NewFirstSpan(size int) *Span {
	return &Span{
		Pos:  new(Cursor),
		size: size,
	}
}

// NewSpan returns a new span of the defined size and in position defined by the cursor value.
func NewSpan(cursor string, size int) (*Span, error) {
	cur, err := ParseCursor(cursor)
	if err != nil {
		return nil, fmt.Errorf("cursor parsing error: %w", err)
	}
	return &Span{
		Pos:  cur,
		size: size,
	}, nil
}

// LimitSize limits size of the span by the maxSize value and returns the size.
func (p *Span) LimitSize(maxSize int) int {
	if maxSize <= 0 {
		panic("maxSize must be positive")
	}
	if maxSize < p.size {
		p.size = maxSize
	}
	if p.size == 0 {
		p.size = maxSize
	}
	return p.size
}

// LimitSize64 limits size of the span by the maxSize value and returns the size.
func (p *Span) LimitSize64(maxSize int64) int64 {
	return int64(p.LimitSize(int(maxSize)))
}

// IsLast returns true in case the span is last in the collection.
func (p *Span) IsLast() bool {
	return p.Pos.IsOutOfScope()
}

// SetLast maks the span as the last in case f is true.
func (p *Span) SetLast(f bool) {
	if f {
		p.Pos.MarkOutOfScope()
	}
}
