package api

import (
	"encoding/xml"
	"fmt"

	"github.com/valyala/fasthttp"
)

type ListBucketsInput struct {
	BucketRegion      *string
	ContinuationToken *string
	MaxBuckets        *string
	Prefix            *string
}

func (input *ListBucketsInput) MarshalHTTP(req *fasthttp.Request) error {
	req.Header.SetMethod(fasthttp.MethodGet)

	args := req.URI().QueryArgs()
	setQuery(args, QueryBucketRegion, input.BucketRegion)
	setQuery(args, QueryContinuationToken, input.ContinuationToken)
	setQuery(args, QueryMaxBuckets, input.MaxBuckets)
	setQuery(args, QueryPrefix, input.Prefix)

	return nil
}

type ListBucketsOutput struct {
	Payload *ListAllMyBucketsResult
}

func (output *ListBucketsOutput) UnmarshalHTTP(resp *fasthttp.Response) error {
	if got, want := resp.StatusCode(), fasthttp.StatusOK; got != want {
		return NewServerSideError(resp)
	}

	var payload ListAllMyBucketsResult
	if err := xml.Unmarshal(resp.Body(), &payload); err != nil {
		return fmt.Errorf("cannot parse ListBucketsOutput body: %w", err)
	}
	output.Payload = &payload
	return nil
}

type ListAllMyBucketsResult struct {
	Buckets           []Bucket `xml:">Bucket"`
	Owner             *Owner
	ContinuationToken *string
	Prefix            *string
}

type Bucket struct {
	BucketRegion *string
	CreationDate *string
	Name         *string
}
