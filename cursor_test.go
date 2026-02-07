package pagination

import "testing"

func TestCursor(t *testing.T) {
	suite := []struct {
		name string
		fn   func(*testing.T, *Cursor) int
		bad  bool
	}{
		{name: "bind", fn: testFetchBind},
		{name: "getset", fn: testFetchGetSet},
		{name: "bind", fn: testFetchBindBad, bad: true},
		{name: "getset", fn: testFetchGetSetBad, bad: true},
	}

	t.Run("short scenario", func(t *testing.T) {
		for _, tc := range suite {
			t.Run(tc.name, func(t *testing.T) {
				if tc.bad {
					t.Run("bad", func(t *testing.T) {
						cur := new(Cursor)
						tc.fn(t, cur)
						defer func() {
							if e := recover(); e == nil {
								t.Error("should panic on improper using")
							}
						}()
						tc.fn(t, cur)
					})
					return
				}
				t.Run("ok", func(t *testing.T) {
					cur := new(Cursor)
					if tc.fn(t, cur) != 0 {
						t.Errorf("should not rewrite default state for the starting position")
					}
					if tc.fn(t, cur) != 1 {
						t.Fatal("should keep state between calls")
					}
				})
			})
		}
	})

	t.Run("long scenario", func(t *testing.T) {
		for _, tc := range suite {
			t.Run(tc.name, func(t *testing.T) {
				if tc.bad {
					t.Run("bad for first span", func(t *testing.T) {
						cur := new(Cursor)
						tc.fn(t, cur)
						defer func() {
							if e := recover(); e == nil {
								t.Error("should panic on improper using")
							}
						}()
						cur.EncodeToString()
					})
					t.Run("bad for some cursor", func(t *testing.T) {
						cur := new(Cursor)
						testFetchBind(t, cur)
						cur, _ = ParseCursor(cur.EncodeToString())
						tc.fn(t, cur)
						defer func() {
							if e := recover(); e == nil {
								t.Error("should panic on improper using")
							}
						}()
						cur.EncodeToString()
					})
					return
				}
				t.Run("ok", func(t *testing.T) {
					cur := new(Cursor)
					tc.fn(t, cur)
					str := cur.EncodeToString()
					newcur, err := ParseCursor(str)
					if err != nil {
						t.Fatal("should not return error:", err)
					}
					if tc.fn(t, newcur) != 1 {
						t.Fatal("should translate state between processes")
					}
				})
			})
		}
	})
}

func testFetchBind(t *testing.T, cur *Cursor) int {
	t.Helper()
	pos := -1
	if err := cur.Bind(&pos); err != nil {
		t.Fatal(err, "should not return error")
	}
	// advance cursor
	pos += 1
	return pos
}

func testFetchBindBad(t *testing.T, cur *Cursor) int {
	t.Helper()
	pos := -1
	if err := cur.Bind(&pos); err != nil {
		t.Fatal(err, "should not return error")
	}
	return pos
}

func testFetchGetSet(t *testing.T, cur *Cursor) int {
	t.Helper()
	pos := -1
	if err := cur.Get(&pos); err != nil {
		t.Fatal(err, "should not return error")
	}
	// advance cursor
	pos += 1
	cur.Set(pos)
	return pos
}

func testFetchGetSetBad(t *testing.T, cur *Cursor) int {
	t.Helper()
	pos := -1
	if err := cur.Get(&pos); err != nil {
		t.Fatal(err, "should not return error")
	}
	return pos
}
