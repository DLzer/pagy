package pagy_test

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/DLzer/pagy"
	"github.com/cornelk/hashmap/assert"
)

var (
	pqPage     = "1"
	pqSize     = "10"
	pqOrderBy  = "first_name"
	pqOrderDir = "desc"
)

type TestUser struct {
	FirstName string
	LastName  string
}

func TestPaginationResponse(t *testing.T) {
	t.Parallel()

	pgres := pagy.PaginationResponse[TestUser]{
		TotalCount: 150,
		TotalPages: 15,
		Page:       1,
		Size:       10,
		HasMore:    true,
		Values: []TestUser{
			{FirstName: "John", LastName: "Wick"},
		},
	}

	assert.Equal(t, pgres.TotalCount, 150)
	assert.Equal(t, pgres.TotalPages, 15)
	assert.Equal(t, pgres.Page, 1)
	assert.Equal(t, pgres.Size, 10)
	assert.Equal(t, pgres.HasMore, true)
	assert.Equal(t, pgres.Values[0].FirstName, "John")
	assert.Equal(t, pgres.Values[0].LastName, "Wick")
}

func TestPaginationQuery(t *testing.T) {
	t.Parallel()

	pq := pagy.PaginationQuery{
		Size:     10,
		Page:     1,
		OrderBy:  pqOrderBy,
		OrderDir: pqOrderDir,
	}

	assert.Equal(t, pq.Size, 10)
	assert.Equal(t, pq.Page, 1)
	assert.Equal(t, pq.OrderBy, pqOrderBy)
	assert.Equal(t, pq.OrderDir, pqOrderDir)
}

func TestPaginatedResponse(t *testing.T) {
	t.Parallel()

	testUsers := []TestUser{
		{FirstName: "John", LastName: "Wick"},
	}

	pq := pagy.PaginationQuery{
		Size:     10,
		Page:     1,
		OrderBy:  pqOrderBy,
		OrderDir: pqOrderDir,
	}

	pgresb := pagy.PaginatedResponse(pq.Size, &pq, testUsers)
	assert.Equal(t, pgresb.Size, 10)
	assert.Equal(t, pgresb.Page, 1)
	assert.Equal(t, pgresb.Values[0].FirstName, "John")
	assert.Equal(t, pgresb.Values[0].LastName, "Wick")
}

func TestDefaultPaginationResponse(t *testing.T) {
	pq := pagy.PaginationQuery{
		Size:     10,
		Page:     1,
		OrderBy:  pqOrderBy,
		OrderDir: pqOrderDir,
	}

	dpr := pagy.DefaultPaginationResponse[TestUser](&pq)

	assert.Equal(t, dpr.TotalCount, 0)
	assert.Equal(t, dpr.TotalPages, 0)
	assert.Equal(t, dpr.Values, []TestUser{})
}

func TestGetPaginationFromRequest(t *testing.T) {
	t.Parallel()

	req, _ := http.NewRequest("GET", "http://localhost", http.NoBody)

	q := req.URL.Query()
	q.Add("page", pqPage)
	q.Add("size", pqSize)
	q.Add("orderBy", pqOrderBy)
	q.Add("orderDir", pqOrderDir)

	req.URL.RawQuery = q.Encode()

	type args struct {
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    *pagy.PaginationQuery
		wantErr bool
	}{
		{
			name:    "base",
			args:    args{req},
			want:    &pagy.PaginationQuery{Page: 1, Size: 10, OrderBy: "first_name DESC"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pagy.GetPaginationFromRequest(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPaginationFromRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPaginationFromRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
