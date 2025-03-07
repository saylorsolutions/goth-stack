package urlprefix

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	defer setReset(t, " / test / ")()
	assert.Equal(t, "/test", Get())
}

func TestApply(t *testing.T) {
	defer setReset(t, "/test")()
	assert.Equal(t, "/test/something", Apply("\tsomething "))
	assert.Equal(t, "/test/something?err=blah blah", Apply("\tsomething?err=blah blah "))
}

func TestGroup(t *testing.T) {
	tests := map[string]struct {
		prefix          string
		requestPath     string
		expectedHandles int
	}{
		"With prefix": {
			prefix:          "/prefix",
			requestPath:     "/prefix/test",
			expectedHandles: 1,
		},
		"Without prefix": {
			prefix:          "",
			requestPath:     "/test",
			expectedHandles: 1,
		},
		"Invalid prefix": {
			prefix:          "\t \r",
			requestPath:     "/prefix/test",
			expectedHandles: 0,
		},
		"Invalid prefix, corrected path": {
			prefix:          "\t \r",
			requestPath:     "/test",
			expectedHandles: 1,
		},
		"Out of scope request": {
			prefix:          "/prefix",
			requestPath:     "/other",
			expectedHandles: 0,
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			defer setReset(t, tc.prefix)()
			var handled int
			mux := http.NewServeMux()
			mux.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
				handled++
			})
			srv := httptest.NewServer(Group(mux))
			defer srv.Close()

			{
				resp, err := http.Get(srv.URL + tc.requestPath)
				assert.NoError(t, err)
				if tc.expectedHandles == 0 {
					assert.Equal(t, 404, resp.StatusCode)
				} else {
					assert.Equal(t, 200, resp.StatusCode)
				}
			}
			assert.Equal(t, tc.expectedHandles, handled)
		})
	}
}

func setReset(t *testing.T, val string) func() {
	orig, isSet := os.LookupEnv(EnvUrlPrefix)
	assert.NoError(t, os.Setenv(EnvUrlPrefix, val))
	initFunc()
	return func() {
		defer initFunc()
		if isSet {
			assert.NoError(t, os.Setenv(EnvUrlPrefix, orig))
			return
		}
		assert.NoError(t, os.Unsetenv(EnvUrlPrefix))
	}
}
