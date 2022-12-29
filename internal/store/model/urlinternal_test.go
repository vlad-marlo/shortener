package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestURL_ShortURL(t *testing.T) {
	tt := []struct {
		name string
		u    *URL
		err  error
	}{
		{
			name: "positive case #1",
			u: &URL{
				BaseURL: "https://marlo.ru",
				User:    "marlo",
			},
			err: nil,
		},
		{
			name: "positive case #2",
			u: &URL{
				BaseURL: "https://ya.ru",
				User:    "marlo",
			},
			err: nil,
		},
		{
			name: "negative case #1",
			u: &URL{
				BaseURL: "https:// ya.ru",
				User:    "marlo",
			},
			err: ErrURLContainSpace,
		},
		{
			name: "negative case #1",
			u: &URL{
				BaseURL: "htt",
				User:    "marlo",
			},
			err: ErrURLTooShort,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.u.ShortURL()
			if tc.err == nil {
				assert.NotEmpty(t, tc.u.ID)
			}
			assert.ErrorIs(t, err, tc.err)
		})
	}
}

func TestNewURL(t *testing.T) {
	tt := []struct {
		name  string
		url   string
		user  string
		corel []string
		err   error
	}{
		{
			name:  "positive case #1",
			url:   "https://example.org",
			user:  "marlo",
			corel: []string{},
			err:   nil,
		},
		{
			name:  "positive case #2",
			url:   "https://example.org",
			user:  "marlo",
			corel: []string{"a"},
			err:   nil,
		},
		{
			name:  "negative case #1",
			url:   "https://example.org",
			user:  "marlo",
			corel: []string{"a", "b"},
			err:   ErrURLBadCorrelationID,
		},
		{
			name:  "negative case #2",
			url:   "htt",
			user:  "marlo",
			corel: []string{},
			err:   ErrURLTooShort,
		},
		{
			name:  "negative case #3",
			url:   "htt ps://example.org",
			user:  "marlo",
			corel: []string{},
			err:   ErrURLContainSpace,
		},
		{
			name:  "negative case #4",
			url:   "htt ps://example.org",
			user:  "marlo",
			corel: []string{"a", "b"},
			err:   ErrURLBadCorrelationID,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			u, err := NewURL(tc.url, tc.user, tc.corel...)
			require.ErrorIs(t, err, tc.err, "got unexpected error")
			if err != nil {
				return
			}
			assert.Equal(t, tc.url, u.BaseURL, "got unexpected base url")
			assert.Equal(t, tc.user, u.User, "got unknown user")
			if len(tc.corel) != 0 {
				assert.Contains(t, tc.corel, u.CorelID, "got bad corel id")
			}
			assert.False(t, u.IsDeleted, "got bad is deleted status")
			assert.NotEmpty(t, u.ID, "shorten id must be not empty")
		})
	}
}
