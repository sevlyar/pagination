package pagination

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash/maphash"
	"reflect"
)

// ErrInvalidCursor is returned when the provided cursor is invalid.
var ErrInvalidCursor = fmt.Errorf("invalid cursor")

// Cursor keeps position of a collection between calls of repository method
// and hides the repository implementation details.
// Cursor allow transfer state across process boundaries.
// Zero value means starting position.
type Cursor struct {
	b []byte

	v        any
	lasthash uint64

	outOfScope bool
}

// ParseCursor parses string representation of the Cursor.
// Empty string treats as starting position.
//
// Returns en errors in case invalid cursor.
func ParseCursor(str string) (*Cursor, error) {
	if str == "" {
		return new(Cursor), nil
	}
	b, err := base64.RawURLEncoding.DecodeString(str)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor: %w", err)
	}
	return &Cursor{b: b}, nil
}

// EncodeToString encodes cursor state to string.
// The string url-encoded and safe to use in url address.
// Empty string means cursor reach end of collection.
//
// EncodeToString panics in case cursor misuse -
// if cursor state was not changed between iterations.
func (cur *Cursor) EncodeToString() string {
	if cur.outOfScope {
		return ""
	}
	if cur.v == nil {
		panic("misuse: cursor should be bind or value set")
	}
	cur.checkValueChanged()
	cur.marshal()
	return base64.RawURLEncoding.EncodeToString(cur.b)
}

// MarkOutOfScope sets flag the cursor reach end of collection.
func (cur *Cursor) MarkOutOfScope() {
	cur.outOfScope = true
}

// IsOutOfScope returns true in case the cursor reach end of collection.
func (cur *Cursor) IsOutOfScope() bool {
	return cur.outOfScope
}

// Get puts state of the cursor into variable pointed by v.
// v should be pointer of comparable value: struct, scalar.
//
// Get panics in case cursor misuse -
// if cursor state was not changed on previous iteration.
//
// Get returns ErrInvalidCursor if the cursor cannot be unmarshalled into v.
//
// Don't use with Bind method.
func (cur *Cursor) Get(v any) error {
	return cur.Bind(v)
}

// Set writes state of the cursor from v.
// v should be comparable value: struct, scalar.
//
// Don't use with Bind method.
func (cur *Cursor) Set(v any) {
	cur.v = v
}

// Bind binds state of the cursor with variable pointed by v.
// The variable value defines internal structure of the cursor.
// Changes of the variable will be change cursor state.
// v should be pointer of comparable value: struct, scalar.
//
// Bind panics in case cursor misuse -
// if cursor state was not changed on previous iteration.
//
// Bind returns ErrInvalidCursor if the cursor cannot be unmarshalled into v.
//
// Don't use with Get/Set methods.
func (cur *Cursor) Bind(v any) error {
	mustBePtr(v)
	cur.checkValueChanged()
	defer cur.updateHash()
	if cur.v != nil {
		ov := reflect.Indirect(reflect.ValueOf(cur.v))
		vv := reflect.ValueOf(v).Elem()
		vv.Set(ov)
		cur.v = v
		return nil
	}
	cur.v = v
	if len(cur.b) == 0 {
		return nil
	}
	if err := json.Unmarshal(cur.b, cur.v); err != nil {
		return ErrInvalidCursor
	}
	return nil
}

func (cur *Cursor) checkValueChanged() {
	if cur.v != nil {
		actual := maphash.Comparable(mapHashSeed, deref(cur.v))
		if cur.lasthash == actual {
			panic("misuse: cursor went unchanged in previous iteration")
		}
	}
}

func (cur *Cursor) updateHash() {
	cur.lasthash = maphash.Comparable(mapHashSeed, deref(cur.v))
}

var mapHashSeed = maphash.MakeSeed()

func mustBePtr(v any) {
	if reflect.ValueOf(v).Kind() != reflect.Pointer {
		panic("value must be a pointer")
	}
}

func (cur *Cursor) marshal() {
	b := bytes.NewBuffer(cur.b[:0])
	if err := json.NewEncoder(b).Encode(cur.v); err != nil {
		panic(err)
	}
	cur.b = b.Bytes()
}

func deref(v any) any {
	if v == nil {
		return nil
	}
	return reflect.Indirect(reflect.ValueOf(v)).Interface()
}
