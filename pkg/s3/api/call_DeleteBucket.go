package api

import (
	"github.com/valyala/fasthttp"
)

type DeleteBucketInput struct {
	// bucket is required
	Bucket string

	ExpectedBucketOwner *string
}

func (input *DeleteBucketInput) GetBucket() string {
	return input.Bucket
}

func (input *DeleteBucketInput) MarshalHTTP(req *fasthttp.Request) error {
	req.Header.SetMethod(fasthttp.MethodDelete)

	setHeader(&req.Header, HeaderXAmzExpectedBucketOwner, input.ExpectedBucketOwner)

	return nil
}

type DeleteBucketOutput struct {
}

func (*DeleteBucketOutput) UnmarshalHTTP(resp *fasthttp.Response) error {
	if got, want := resp.StatusCode(), fasthttp.StatusNoContent; got != want {
		return NewServerSideError(resp)
	}

	return nil
}
