package api

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAPIError(t *testing.T) {
	testCases := []struct {
		input APIError
		want  string
	}{
		{
			input: APIError{},
			want:  "s3 api error",
		},
		{
			input: APIError{
				Message: "my-message",
			},
			want: "s3 api error: my-message",
		},
		{
			input: APIError{
				Code:      "my-code",
				RequestID: "my-request-id",
				HostID:    "my-host-id",
			},
			want: "s3 api error (ErrorCode:my-code) (RequestID:my-request-id) (HostID:my-host-id)",
		},
		{
			input: APIError{
				Code:      "my-code",
				Message:   "my-message",
				RequestID: "my-request-id",
				HostID:    "my-host-id",
			},
			want: "s3 api error (ErrorCode:my-code) (RequestID:my-request-id) (HostID:my-host-id): my-message",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			require.EqualError(t, &tc.input, tc.want)
		})
	}
}
