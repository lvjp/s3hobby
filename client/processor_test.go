package client

import (
	"encoding/xml"
	"testing"

	"github.com/lvjp/s3hobby/pkg/s3/api"
	"github.com/lvjp/s3hobby/pkg/s3/signing"
	"github.com/lvjp/s3hobby/pkg/s3/signing/anonymous"

	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

var _ api.HTTPRequestMarshaler = (*dummyInput)(nil)

type dummyInput struct{}

func (d *dummyInput) MarshalHTTP(req *fasthttp.Request) error {
	return nil
}

var _ api.HTTPResponseUnmarshaler = (*dummyOutput)(nil)

type dummyOutput struct{}

func (d *dummyOutput) UnmarshalHTTP(resp *fasthttp.Response) error {
	if resp.StatusCode() != fasthttp.StatusOK {
		return api.NewServerSideError(resp)
	}
	return nil
}

func TestPerformCall(t *testing.T) {
	baseOptions := &Options{
		EndpointHost:   "127.0.0.1:1",
		UsePathStyle:   true,
		SiginingRegion: "eu-west-1",
		Signer:         anonymous.NewSigner(),
		Credentials:    &signing.Credentials{},
	}

	baseOptions.setDefaults()

	require.NoError(t, baseOptions.Validate())

	t.Run("HTTPError", func(t *testing.T) {
		expectedError := &api.ServerSideError{
			APIError: &api.APIError{
				Code:      "xml-code",
				Message:   "xml-message",
				RequestID: "xml-request-id",
				HostID:    "xml-host-id",
			},

			HTTPStatusCode: fasthttp.StatusServiceUnavailable,
			XAmzRequestID:  "header-request-id",
			XAmzID2:        "header-amz-id2",
		}

		respBody, err := xml.Marshal(expectedError.APIError)
		require.NoError(t, err)

		opts := baseOptions.With(func(o *Options) {
			o.HTTPClient = HTTPClientFunc(func(req *fasthttp.Request, resp *fasthttp.Response) error {
				resp.SetStatusCode(expectedError.HTTPStatusCode)
				resp.Header.Set(api.HeaderXAmzRequestID, expectedError.XAmzRequestID)
				resp.Header.Set(api.HeaderXAmzID2, expectedError.XAmzID2)
				resp.SetBody(respBody)
				return nil
			})
		})

		var output *dummyOutput
		var metadata *Metadata
		output, metadata, err = PerformCall[*dummyInput, *dummyOutput](t.Context(), opts, &dummyInput{})
		require.Error(t, err)
		require.Equal(t, expectedError, err)
		t.Logf("Output: %+v", output)
		t.Logf("Metadata: %+v", metadata)

		require.Nil(t, output)
		require.NotNil(t, metadata)
	})
}
