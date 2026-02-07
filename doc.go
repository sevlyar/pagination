/*
Package pagination allows construct the simple and safe cursor-based pagination for services.

Pagination parameter should be the last in repository method parameters:

	type ItemFilterer interface {
		FilterItems(ctx context.Context, params FilterParams, span *pagination.Span) ([]Items, error)
	}

Basic repository implementation:

	func (s *Storage) FilterItems(ctx context.Context, params FilterParams, span *pagination.Span) ([]Items, error) {
		var id int64
		if err := span.Pos.Bind(&id); err != nil {
			return nil, ErrInvalidCursor
		}

		const spanSizeLimit = 100
		limit := span.LimitSize(spanSizeLimit)

		// use id and limit in database query, fetch items
		// ...

		if len(items) > 0 {
			// advance cursor
			id = items[len(items)-1].Id
		}

		span.SetLast(len(items) < limit)

		return items, nil
	}

Repository method using:

	func handleGetItems(resp http.ResponseWriter, req *http.Request) {
		// GET /items?limit=10&cursor=dGltZXN0YW1wOjE2MjI1NDg4MDA=

		q := req.URL.Query()
		cursor := q.Get("cursor")
		limitStr := q.Get("limit")
		limit, _ := strconv.Atoi(limitStr)

		span, _ := NewSpan(cursor, limit)
		items, _ := stor.FilterItems(req.Context(), FilterParams{}, span)

		resp.Header().Add("X-Pagination-Cursor", span.Pos.EncodeToString())

		// marshal response data
		// ...
	}

Compound cursor allows to use any sorting fields and order:

	type Cursor struct {
		T int64 // time
		R int64 // rowid
	}
	cursor := Cursor {
		T: time.Now().UnixMilli(),
		R: math.MaxInt64,
	}
	span.Pos.Bind(&cursor)

The package supports only comparable values as cursor state. For example,
unable to use slice as cursor or it's part. But you can use arrays.
*/
package pagination
