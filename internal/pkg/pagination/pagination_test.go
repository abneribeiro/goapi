package pagination

import (
	"net/http"
	"net/url"
	"testing"
)

func TestFromRequest(t *testing.T) {
	tests := []struct {
		name            string
		queryParams     map[string]string
		expectedPage    int
		expectedPerPage int
		expectedOffset  int
	}{
		{
			name:            "default values",
			queryParams:     map[string]string{},
			expectedPage:    1,
			expectedPerPage: 10,
			expectedOffset:  0,
		},
		{
			name:            "custom page and per_page",
			queryParams:     map[string]string{"page": "3", "per_page": "20"},
			expectedPage:    3,
			expectedPerPage: 20,
			expectedOffset:  40,
		},
		{
			name:            "negative page defaults to 1",
			queryParams:     map[string]string{"page": "-1"},
			expectedPage:    1,
			expectedPerPage: 10,
			expectedOffset:  0,
		},
		{
			name:            "per_page exceeds max",
			queryParams:     map[string]string{"per_page": "500"},
			expectedPage:    1,
			expectedPerPage: 100,
			expectedOffset:  0,
		},
		{
			name:            "invalid page string",
			queryParams:     map[string]string{"page": "abc"},
			expectedPage:    1,
			expectedPerPage: 10,
			expectedOffset:  0,
		},
		{
			name:            "zero per_page defaults",
			queryParams:     map[string]string{"per_page": "0"},
			expectedPage:    1,
			expectedPerPage: 10,
			expectedOffset:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := url.Values{}
			for k, v := range tt.queryParams {
				values.Set(k, v)
			}

			req := &http.Request{
				URL: &url.URL{
					RawQuery: values.Encode(),
				},
			}

			params := FromRequest(req)

			if params.Page != tt.expectedPage {
				t.Errorf("expected page %d, got %d", tt.expectedPage, params.Page)
			}
			if params.PerPage != tt.expectedPerPage {
				t.Errorf("expected per_page %d, got %d", tt.expectedPerPage, params.PerPage)
			}
			if params.Offset != tt.expectedOffset {
				t.Errorf("expected offset %d, got %d", tt.expectedOffset, params.Offset)
			}
		})
	}
}

func TestCalculateTotalPages(t *testing.T) {
	tests := []struct {
		name     string
		total    int64
		perPage  int
		expected int
	}{
		{"zero total", 0, 10, 0},
		{"exact division", 100, 10, 10},
		{"with remainder", 105, 10, 11},
		{"less than per_page", 5, 10, 1},
		{"single item", 1, 10, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateTotalPages(tt.total, tt.perPage)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}
