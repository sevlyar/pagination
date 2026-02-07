# pagination [![Go Reference](https://pkg.go.dev/badge/github.com/sevlyar/pagination.svg)](https://pkg.go.dev/github.com/sevlyar/pagination)

The pagination package allows construct the simple and safe cursor-based pagination for services.

Cursor-based pagination is a highly efficient, high-performance technique for large datasets
that uses a pointer (the "cursor") to mark the last retrieved record. Unlike offset pagination, 
it avoids expensive, slow database scans on deep pages, providing consistent performance by directly 
querying records after a specific identifier, such as a timestamp or ID.

### How It Works

Instead of using page numbers, the client requests a specific number of items after a unique marker (the cursor).

* **Initial Request**: GET /items?limit=10
* **Subsequent Request**: GET /items?limit=10&cursor=dGltZXN0YW1wOjE2MjI1NDg4MDA=
* **Response**: The server returns the next 10 items and a next cursor value

## Installation

```sh
go get github.com/sevlyar/pagination
```

## Documentation

Full `go doc` style documentation for the package can be viewed online by using the GoDoc site:
https://pkg.go.dev/github.com/sevlyar/pagination