package api

import (
	"encoding/xml"

	"github.com/valyala/fasthttp"
)

type CreateBucketInput struct {
	// Bucket is required
	Bucket string

	ACL                        *string
	GrantFullControl           *string
	GrantRead                  *string
	GrantReadACP               *string
	GrantWrite                 *string
	GrantWriteACP              *string
	ObjectLockEnabledForBucket *string
	ObjectOwnership            *ObjectOwnership

	CreateBucketConfiguration *struct {
		LocationConstraint *LocationConstraint
	}
}

func (input *CreateBucketInput) GetBucket() string {
	return input.Bucket
}

func (input *CreateBucketInput) MarshalHTTP(req *fasthttp.Request) error {
	req.Header.SetMethod(fasthttp.MethodPut)

	setHeader(&req.Header, HeaderXAmzACL, input.ACL)
	setHeader(&req.Header, HeaderXAmzGrantFullControl, input.GrantFullControl)
	setHeader(&req.Header, HeaderXAmzGrantRead, input.GrantRead)
	setHeader(&req.Header, HeaderXAmzGrantReadACP, input.GrantReadACP)
	setHeader(&req.Header, HeaderXAmzGrantWrite, input.GrantWrite)
	setHeader(&req.Header, HeaderXAmzGrantWriteACP, input.GrantWriteACP)
	setHeader(&req.Header, HeaderXAmzBucketObjectLockEnabled, input.ObjectLockEnabledForBucket)
	setHeader(&req.Header, HeaderXAmzObjectOwnership, (*string)(input.ObjectOwnership))

	if input.CreateBucketConfiguration != nil {
		inputBody, err := xml.Marshal(input.CreateBucketConfiguration)
		if err != nil {
			return err
		}

		req.SetBody(inputBody)
	}

	return nil
}

type CreateBucketOutput struct {
	Location *string
}

func (output *CreateBucketOutput) UnmarshalHTTP(resp *fasthttp.Response) error {
	if got, want := resp.StatusCode(), fasthttp.StatusOK; got != want {
		return NewServerSideError(resp)
	}

	output.Location = extractHeader(&resp.Header, HeaderLocation)

	return nil
}
