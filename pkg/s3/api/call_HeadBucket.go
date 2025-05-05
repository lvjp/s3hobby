package api

import (
	"github.com/valyala/fasthttp"
)

type HeadBucketInput struct {
	// Bucket is mandatory
	Bucket string

	ExpectedBucketOwner *string
}

func (input *HeadBucketInput) GetBucket() string {
	return input.Bucket
}

func (input *HeadBucketInput) MarshalHTTP(req *fasthttp.Request) error {
	req.Header.SetMethod(fasthttp.MethodHead)

	setHeader(&req.Header, HeaderXAmzExpectedBucketOwner, input.ExpectedBucketOwner)

	return nil
}

type HeadBucketOutput struct {
	AccessPointAlias *string
	BucketRegion     *string
}

func (output *HeadBucketOutput) UnmarshalHTTP(resp *fasthttp.Response) error {
	if got, want := resp.StatusCode(), fasthttp.StatusOK; got != want {
		return NewServerSideError(resp)
	}

	output.AccessPointAlias = extractHeader(&resp.Header, HeaderXAmzAccessPointAlias)
	output.BucketRegion = extractHeader(&resp.Header, HeaderXAmzBucketRegion)

	return nil
}
