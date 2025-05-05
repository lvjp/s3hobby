package api

import (
	"github.com/valyala/fasthttp"
)

type PutObjectInput struct {
	// Bucket is mandatory
	Bucket string

	// Key is mandatory
	Key string

	Body []byte

	CacheControl       *string
	ContentDisposition *string
	ContentEncoding    *string
	ContentLanguage    *string
	ContentMD5         *string
	ContentType        *string
	Expires            *string
	IfMatch            *string
	IfNoneMatch        *string

	ACL                       *string
	ChecksumAlgorithm         *string
	ChecksumCRC32             *string
	ChecksumCRC32C            *string
	ChecksumCRC64NVME         *string
	ChecksumSHA1              *string
	ChecksumSHA256            *string
	ExpectedBucketOwner       *string
	GrantFullControl          *string
	GrantRead                 *string
	GrantReadACP              *string
	GrantWriteACP             *string
	ObjectLockLegalHold       *ObjectLockLegalHoldStatus
	ObjectLockMode            *ObjectLockMode
	ObjectLockRetainUntilDate *string
	RequestPayer              *string
	ServerSideEncryption      *SSEType
	SSEBucketKeyEnabled       *string
	SSECustomerAlgorithm      *SSECustomerAlgorithm
	SSECustomerKey            *string
	SSECustomerKeyMD5         *string
	SSEKMSEncryptionContext   *string
	SSEKMSKeyId               *string
	StorageClass              *StorageClass
	Tagging                   *string
	WebsiteRedirectLocation   *string
	WriteOffsetBytes          *string

	TrailerChecksumCRC32     *string
	TrailerChecksumCRC64NVME *string
	TrailerChecksumSHA1      *string
	TrailerChecksumCRC32C    *string
	TrailerChecksumSHA256    *string
}

func (input *PutObjectInput) GetBucket() string {
	return input.Bucket
}

func (input *PutObjectInput) GetKey() string {
	return input.Key
}

func (input *PutObjectInput) MarshalHTTP(req *fasthttp.Request) error {
	req.Header.SetMethod(fasthttp.MethodPut)

	if input.Body != nil {
		req.SetBody(input.Body)
	} else {
		// Ensure that the Content-Length Header is set to 0
		req.SetBody([]byte{})
	}

	setHeader(&req.Header, HeaderCacheControl, input.CacheControl)
	setHeader(&req.Header, HeaderContentDisposition, input.ContentDisposition)
	setHeader(&req.Header, HeaderContentEncoding, input.ContentEncoding)
	setHeader(&req.Header, HeaderContentLanguage, input.ContentLanguage)
	setHeader(&req.Header, HeaderContentMD5, input.ContentMD5)
	setHeader(&req.Header, HeaderContentType, input.ContentType)
	setHeader(&req.Header, HeaderExpires, input.Expires)
	setHeader(&req.Header, HeaderIfMatch, input.IfMatch)
	setHeader(&req.Header, HeaderIfNoneMatch, input.IfNoneMatch)

	setHeader(&req.Header, HeaderXAmzACL, input.ACL)
	setHeader(&req.Header, HeaderXAmzSdkChecksumAlgorithm, input.ChecksumAlgorithm)
	setHeaderOrTrailer(&req.Header, HeaderXAmzChecksumCRC32, input.ChecksumCRC32, input.TrailerChecksumCRC32)
	setHeaderOrTrailer(&req.Header, HeaderXAmzChecksumCRC32C, input.ChecksumCRC32C, input.TrailerChecksumCRC32C)
	setHeaderOrTrailer(&req.Header, HeaderXAmzChecksumCRC64NVME, input.ChecksumCRC64NVME, input.TrailerChecksumCRC64NVME)
	setHeaderOrTrailer(&req.Header, HeaderXAmzChecksumSHA1, input.ChecksumSHA1, input.TrailerChecksumSHA1)
	setHeaderOrTrailer(&req.Header, HeaderXAmzChecksumSHA256, input.ChecksumSHA256, input.TrailerChecksumSHA256)
	setHeader(&req.Header, HeaderXAmzExpectedBucketOwner, input.ExpectedBucketOwner)
	setHeader(&req.Header, HeaderXAmzGrantFullControl, input.GrantFullControl)
	setHeader(&req.Header, HeaderXAmzGrantRead, input.GrantRead)
	setHeader(&req.Header, HeaderXAmzGrantReadACP, input.GrantReadACP)
	setHeader(&req.Header, HeaderXAmzGrantWriteACP, input.GrantWriteACP)
	setHeader(&req.Header, HeaderXAmzObjectLockLegalHold, (*string)(input.ObjectLockLegalHold))
	setHeader(&req.Header, HeaderXAmzObjectLockMode, (*string)(input.ObjectLockMode))
	setHeader(&req.Header, HeaderXAmzObjectLockRetainUntilDate, input.ObjectLockRetainUntilDate)
	setHeader(&req.Header, HeaderXAmzRequestPayer, input.RequestPayer)
	setHeader(&req.Header, HeaderXAmzServerSideEncryption, (*string)(input.ServerSideEncryption))
	setHeader(&req.Header, HeaderXAmzSSEBucketKeyEnabled, input.SSEBucketKeyEnabled)
	setHeader(&req.Header, HeaderXAmzSSECustomerAlgorithm, (*string)(input.SSECustomerAlgorithm))
	setHeader(&req.Header, HeaderXAmzSSECustomerKey, input.SSECustomerKey)
	setHeader(&req.Header, HeaderXAmzSSECustomerKeyMD5, input.SSECustomerKeyMD5)
	setHeader(&req.Header, HeaderXAmzSSEKMSEncryptionContext, input.SSEKMSEncryptionContext)
	setHeader(&req.Header, HeaderXAmzSSEKMSKeyId, input.SSEKMSKeyId)
	setHeader(&req.Header, HeaderXAmzStorageClass, (*string)(input.StorageClass))
	setHeader(&req.Header, HeaderXAmzTagging, input.Tagging)
	setHeader(&req.Header, HeaderXAmzWebsiteRedirectLocation, input.WebsiteRedirectLocation)
	setHeader(&req.Header, HeaderXAmzWriteOffsetBytes, input.WriteOffsetBytes)

	return nil
}

type PutObjectOutput struct {
	ETag *string

	ChecksumCRC32           *string
	ChecksumCRC32C          *string
	ChecksumCRC64NVME       *string
	ChecksumSHA1            *string
	ChecksumSHA256          *string
	ChecksumType            *ChecksumType
	Expiration              *string
	ObjectSize              *string
	RequestPayer            *RequestPayer
	SSE                     *SSEType
	SSEBucketKeyEnabled     *string
	SSECustomerAlgorithm    *SSECustomerAlgorithm
	SSECustomerKeyMD5       *string
	SSEKMSEncryptionContext *string
	SSEKMSKeyId             *string
	VersionId               *string
}

func (output *PutObjectOutput) UnmarshalHTTP(resp *fasthttp.Response) error {
	if got, want := resp.StatusCode(), fasthttp.StatusOK; got != want {
		return NewServerSideError(resp)
	}

	output.ETag = extractHeader(&resp.Header, HeaderETag)

	output.ChecksumCRC32 = extractHeader(&resp.Header, HeaderXAmzChecksumCRC32)
	output.ChecksumCRC32C = extractHeader(&resp.Header, HeaderXAmzChecksumCRC32C)
	output.ChecksumCRC64NVME = extractHeader(&resp.Header, HeaderXAmzChecksumCRC64NVME)
	output.ChecksumSHA1 = extractHeader(&resp.Header, HeaderXAmzChecksumSHA1)
	output.ChecksumSHA256 = extractHeader(&resp.Header, HeaderXAmzChecksumSHA256)
	output.ChecksumType = (*ChecksumType)(extractHeader(&resp.Header, HeaderXAmzChecksumType))
	output.Expiration = extractHeader(&resp.Header, HeaderXAmzExpiration)
	output.ObjectSize = extractHeader(&resp.Header, HeaderXAmzSize)
	output.RequestPayer = (*RequestPayer)(extractHeader(&resp.Header, HeaderXAmzRequestCharged))
	output.SSE = (*SSEType)(extractHeader(&resp.Header, HeaderXAmzServerSideEncryption))
	output.SSEBucketKeyEnabled = extractHeader(&resp.Header, HeaderXAmzSSEBucketKeyEnabled)
	output.SSECustomerAlgorithm = (*SSECustomerAlgorithm)(extractHeader(&resp.Header, HeaderXAmzSSECustomerAlgorithm))
	output.SSECustomerKeyMD5 = extractHeader(&resp.Header, HeaderXAmzSSECustomerKeyMD5)
	output.SSEKMSEncryptionContext = extractHeader(&resp.Header, HeaderXAmzSSEKMSEncryptionContext)
	output.SSEKMSKeyId = extractHeader(&resp.Header, HeaderXAmzSSEKMSKeyId)
	output.VersionId = extractHeader(&resp.Header, HeaderXAmzVersionId)

	return nil
}
