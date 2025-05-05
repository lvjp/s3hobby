package api

import (
	"encoding/xml"

	"github.com/valyala/fasthttp"
)

type GetBucketLocationInput struct {
	// Bucket is mandatory
	Bucket string

	ExpectedBucketOwner *string
}

func (input *GetBucketLocationInput) GetBucket() string {
	return input.Bucket
}

func (input *GetBucketLocationInput) MarshalHTTP(req *fasthttp.Request) error {
	req.Header.SetMethod(fasthttp.MethodGet)

	args := req.URI().QueryArgs()
	setResource(args, QueryLocation)

	setHeader(&req.Header, HeaderXAmzExpectedBucketOwner, input.ExpectedBucketOwner)

	return nil
}

type GetBucketLocationOutput struct {
	XMLName            xml.Name            `xml:"LocationConstraint"`
	LocationConstraint *LocationConstraint `xml:",chardata"`
}

func (output *GetBucketLocationOutput) UnmarshalHTTP(resp *fasthttp.Response) error {
	if got, want := resp.StatusCode(), fasthttp.StatusOK; got != want {
		return NewServerSideError(resp)
	}

	return xml.Unmarshal(resp.Body(), output)
}
