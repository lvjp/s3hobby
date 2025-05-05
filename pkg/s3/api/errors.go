package api

import (
	"encoding/xml"
	"fmt"

	"github.com/valyala/fasthttp"
)

type APIError struct {
	Code      string `xml:"Code"`
	Message   string `xml:"Message"`
	RequestID string `xml:"RequestId"`
	HostID    string `xml:"HostId"`
}

func (e *APIError) Error() string {
	ret := "s3 api error"

	if e.Code != "" {
		ret += fmt.Sprintf(" (ErrorCode:%s)", e.Code)
	}

	if e.RequestID != "" {
		ret += fmt.Sprintf(" (RequestID:%s)", e.RequestID)
	}

	if e.HostID != "" {
		ret += fmt.Sprintf(" (HostID:%s)", e.HostID)
	}

	if e.Message != "" {
		ret += ": " + e.Message
	}

	return ret
}

type ServerSideError struct {
	APIError *APIError

	HTTPStatusCode int
	XAmzRequestID  string
	XAmzID2        string
}

func NewServerSideError(resp *fasthttp.Response) *ServerSideError {
	statusCode := resp.StatusCode()

	ret := &ServerSideError{
		HTTPStatusCode: statusCode,
		XAmzRequestID:  string(resp.Header.Peek(HeaderXAmzRequestID)),
		XAmzID2:        string(resp.Header.Peek(HeaderXAmzID2)),
	}

	if statusCode >= 200 && statusCode != fasthttp.StatusNoContent {
		ret.APIError = new(APIError)
		if err := xml.Unmarshal(resp.Body(), ret.APIError); err != nil {
			ret.APIError.Message = fmt.Sprintf("xml error response deserializing error: %v", err)
		}
	}

	return ret
}

func (e *ServerSideError) Error() string {
	apiError := "error message not found"
	if e.APIError != nil {
		apiError = e.APIError.Error()
	}

	var statusCode string
	if e.HTTPStatusCode != 0 {
		statusCode = fmt.Sprintf(" (StatusCode:%d)", e.HTTPStatusCode)
	}

	return fmt.Sprintf("server-side error occurred%s: %s", statusCode, apiError)
}

func (e *ServerSideError) Unwrap() error {
	return e.APIError
}
