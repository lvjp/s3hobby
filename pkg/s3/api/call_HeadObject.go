package api

import (
	"github.com/valyala/fasthttp"
)

type HeadObjectInput struct {
	// Bucket is mandatory
	Bucket string

	// Key is mandatory
	Key string

	PartNumber                 *string
	ResponseCacheControl       *string
	ResponseContentDisposition *string
	ResponseContentEncoding    *string
	ResponseContentLanguage    *string
	ResponseContentType        *string
	ResponseExpires            *string
	VersionId                  *string

	IfMatch           *string
	IfModifiedSince   *string
	IfNoneMatch       *string
	IfUnmodifiedSince *string
	Range             *string

	ChecksumMode         *ChecksumMode
	ExpectedBucketOwner  *string
	RequestPayer         *RequestPayer
	SSECustomerAlgorithm *SSECustomerAlgorithm
	SSECustomerKey       *string
	SSECustomerKeyMD5    *string
}

func (input *HeadObjectInput) GetBucket() string {
	return input.Bucket
}

func (input *HeadObjectInput) GetKey() string {
	return input.Key
}

func (input *HeadObjectInput) MarshalHTTP(req *fasthttp.Request) error {
	req.Header.SetMethod(fasthttp.MethodHead)

	args := req.URI().QueryArgs()
	setQuery(args, QueryPartNumber, input.PartNumber)
	setQuery(args, QueryResponseCacheControl, input.ResponseCacheControl)
	setQuery(args, QueryResponseContentDisposition, input.ResponseContentDisposition)
	setQuery(args, QueryResponseContentEncoding, input.ResponseContentEncoding)
	setQuery(args, QueryResponseContentLanguage, input.ResponseContentLanguage)
	setQuery(args, QueryResponseContentType, input.ResponseContentType)
	setQuery(args, QueryResponseExpires, input.ResponseExpires)
	setQuery(args, QueryVersionID, input.VersionId)

	setHeader(&req.Header, HeaderIfMatch, input.IfMatch)
	setHeader(&req.Header, HeaderIfModifiedSince, input.IfModifiedSince)
	setHeader(&req.Header, HeaderIfNoneMatch, input.IfNoneMatch)
	setHeader(&req.Header, HeaderIfUnmodifiedSince, input.IfUnmodifiedSince)
	setHeader(&req.Header, HeaderRange, input.Range)

	setHeader(&req.Header, HeaderXAmzChecksumMode, (*string)(input.ChecksumMode))
	setHeader(&req.Header, HeaderXAmzExpectedBucketOwner, input.ExpectedBucketOwner)
	setHeader(&req.Header, HeaderXAmzRequestPayer, (*string)(input.RequestPayer))
	setHeader(&req.Header, HeaderXAmzSSECustomerAlgorithm, (*string)(input.SSECustomerAlgorithm))
	setHeader(&req.Header, HeaderXAmzSSECustomerKey, input.SSECustomerKey)
	setHeader(&req.Header, HeaderXAmzSSECustomerKeyMD5, input.SSECustomerKeyMD5)

	return nil
}

type HeadObjectOutput struct {
	AcceptRanges       *string
	CacheControl       *string
	ContentDisposition *string
	ContentEncoding    *string
	ContentLanguage    *string
	ContentLength      *string
	ContentRange       *string
	ContentType        *string
	ETag               *string
	Expires            *string
	LastModified       *string

	ArchiveStatus             *string
	ChecksumCRC32             *string
	ChecksumCRC32C            *string
	ChecksumCRC64NVME         *string
	ChecksumSHA1              *string
	ChecksumSHA256            *string
	ChecksumType              *ChecksumType
	DeleteMarker              *string
	Expiration                *string
	MissingMeta               *string
	ObjectLockLegalHoldStatus *ObjectLockLegalHoldStatus
	ObjectLockMode            *ObjectLockMode
	ObjectLockRetainUntilDate *string
	PartsCount                *string
	ReplicationStatus         *ReplicationStatus
	RequestCharged            *RequestPayer
	Restore                   *string
	SSE                       *SSEType
	SSEBucketKeyEnabled       *string
	SSECustomerAlgorithm      *SSECustomerAlgorithm
	SSECustomerKeyMD5         *string
	SSEKMSKeyId               *string
	StorageClass              *StorageClass
	VersionId                 *string
	WebsiteRedirectLocation   *string
}

func (output *HeadObjectOutput) UnmarshalHTTP(resp *fasthttp.Response) error {
	if got, want := resp.StatusCode(), fasthttp.StatusOK; got != want {
		return NewServerSideError(resp)
	}

	output.AcceptRanges = extractHeader(&resp.Header, HeaderAcceptRanges)
	output.CacheControl = extractHeader(&resp.Header, HeaderCacheControl)
	output.ContentDisposition = extractHeader(&resp.Header, HeaderContentDisposition)
	output.ContentEncoding = extractHeader(&resp.Header, HeaderContentEncoding)
	output.ContentLanguage = extractHeader(&resp.Header, HeaderContentLanguage)
	output.ContentLength = extractHeader(&resp.Header, HeaderContentLength)
	output.ContentRange = extractHeader(&resp.Header, HeaderContentRange)
	output.ContentType = extractHeader(&resp.Header, HeaderContentType)
	output.ETag = extractHeader(&resp.Header, HeaderETag)
	output.Expires = extractHeader(&resp.Header, HeaderExpires)
	output.LastModified = extractHeader(&resp.Header, HeaderLastModified)

	output.ArchiveStatus = extractHeader(&resp.Header, HeaderXAmzArchiveStatus)
	output.ChecksumCRC32 = extractHeader(&resp.Header, HeaderXAmzChecksumCRC32)
	output.ChecksumCRC32C = extractHeader(&resp.Header, HeaderXAmzChecksumCRC32C)
	output.ChecksumCRC64NVME = extractHeader(&resp.Header, HeaderXAmzChecksumCRC64NVME)
	output.ChecksumSHA1 = extractHeader(&resp.Header, HeaderXAmzChecksumSHA1)
	output.ChecksumSHA256 = extractHeader(&resp.Header, HeaderXAmzChecksumSHA256)
	output.ChecksumType = (*ChecksumType)(extractHeader(&resp.Header, HeaderXAmzChecksumType))
	output.DeleteMarker = extractHeader(&resp.Header, HeaderXAmzDeleteMarker)
	output.Expiration = extractHeader(&resp.Header, HeaderXAmzExpiration)
	output.MissingMeta = extractHeader(&resp.Header, HeaderXAmzMissingMeta)
	output.ObjectLockLegalHoldStatus = (*ObjectLockLegalHoldStatus)(extractHeader(&resp.Header, HeaderXAmzObjectLockLegalHold))
	output.ObjectLockMode = (*ObjectLockMode)(extractHeader(&resp.Header, HeaderXAmzObjectLockMode))
	output.ObjectLockRetainUntilDate = extractHeader(&resp.Header, HeaderXAmzObjectLockRetainUntilDate)
	output.PartsCount = extractHeader(&resp.Header, HeaderXAmzPartsCount)
	output.ReplicationStatus = (*ReplicationStatus)(extractHeader(&resp.Header, HeaderXAmzReplicationStatus))
	output.RequestCharged = (*RequestPayer)(extractHeader(&resp.Header, HeaderXAmzRequestCharged))
	output.Restore = extractHeader(&resp.Header, HeaderXAmzRestore)
	output.SSE = (*SSEType)(extractHeader(&resp.Header, HeaderXAmzServerSideEncryption))
	output.SSEBucketKeyEnabled = extractHeader(&resp.Header, HeaderXAmzSSEBucketKeyEnabled)
	output.SSECustomerAlgorithm = (*SSECustomerAlgorithm)(extractHeader(&resp.Header, HeaderXAmzSSECustomerAlgorithm))
	output.SSECustomerKeyMD5 = extractHeader(&resp.Header, HeaderXAmzSSECustomerKeyMD5)
	output.SSEKMSKeyId = extractHeader(&resp.Header, HeaderXAmzSSEKMSKeyId)
	output.StorageClass = (*StorageClass)(extractHeader(&resp.Header, HeaderXAmzStorageClass))
	output.VersionId = extractHeader(&resp.Header, HeaderXAmzVersionId)
	output.WebsiteRedirectLocation = extractHeader(&resp.Header, HeaderXAmzWebsiteRedirectLocation)

	return nil
}
