package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantError bool
	}{
		{
			name:      "empty url",
			url:       "",
			wantError: true,
		},
		{
			name:      "length < 4",
			url:       "y.a",
			wantError: true,
		},
		{
			name:      "space in url",
			url:       "yandex. ru",
			wantError: true,
		},
		{
			name:      "positive case #1",
			url:       "y.ru",
			wantError: false,
		},
		{
			name:      "positive case #2",
			url:       "yandex.ru",
			wantError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := URL{
				BaseURL: tt.url,
			}
			err := url.Validate()
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
