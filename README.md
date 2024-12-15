# pagy
[![Test](https://github.com/DLzer/pagy/actions/workflows/test.yml/badge.svg)](https://github.com/DLzer/pagy/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/DLzer/pagy)](https://goreportcard.com/report/github.com/DLzer/pagy)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/DLzer/pagy)](https://pkg.go.dev/github.com/DLzer/pagy)

![golang-gopher](https://static-projects.nyc3.cdn.digitaloceanspaces.com/golang_gopher300x300.png)]


## Overview
Go is idiomatic, opinionated, and frequents a do-it-yourself approach. One case in particular is when it comes to implementing pagination. That's where pagy comes in. It's a strongly opinionated utility package aimed at getting rid of the pesky task of writing your own pagination utilities.

pagy works with the stdlib [http.Request](https://pkg.go.dev/net/http#Request) so it can easily be dropped into any project or handler that implements such.

## Features

* Scans the http.Request for pagination values
* Formats pagination values into simplified structure
* Can be passed into ctx or as a parameter to your repository layer
* Easily extendable

## Install
```bash
go get github.com/DLzer/pagy
```

## What it looks like
pagy will output a consistent top level structure containing important pagination info to your client
```json
{
    "total_count": 150,
    "total_pages": 15,
    "page": 1,
    "size": 10
    "has_more": true
    "values": {
        ...
    }
}
```

## How to use with `http.Request`
```go
package main

import (
    "fmt"
    "net/http"
)

// @Param page query int
// @Param size query int
// @Param orderBy query string
// @Param orderDir query string optional
// @Example GET /users?page=1&size=10&orderBy=first_name
func getUsersList(w http.ResponseWriter, req *http.Request) {
    pq, err := pagy.GetPaginationFromRequest(r)
	if err != nil {
		http.Error(w, err.Error())
		return
	}
}

func main() {
    http.HandleFunc("/users", getUsersList)
    http.ListenAndServe(":8090", nil)
}
```

## How to use in repository layer
pagy uses generics so it's agnostic and doesn't really care what type of structure your data response will be. Instead of responding with nil you can use the `DefaultPaginationResponse` to return a basic empty structure that your client can interpret.
```go
func (p *postgresUserRepository) GetList(ctx context.Context, pq *pagy.PaginationQuery) (*pagy.PaginationResponse[domain.User], error) {{
    var count int
    countQuery := "SELECT count(uuid) FROM users;"
    if err := p.conn.QueryRow(ctx, countQuery).Scan(&count); err != nil {
    	return nil, err
    }

    if count == 0 {
    	return pagy.DefaultPaginationResponse[domain.User](pq), nil
    }

    query := "SELECT * FROM users ORDER BY $1::text OFFSET $2 LIMIT $3"
    rows, err := p.conn.Query(ctx, query, pq.GetOrderBy(), pq.GetOffset(), pq.GetLimit())
    if err != nil {
    	return nil, err
    }
    defer rows.Close()

    var uu []domain.User
    for rows.Next() {
        //... Scan your user
    }

    return pagy.PaginatedResponse(count, pq, uu), nil
}
```

## License

This project is licensed under the terms of the MIT license.