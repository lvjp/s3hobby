package api

type RequestPayer string

const (
	RequestPayerRequester RequestPayer = "requester"
)

type SSECustomerAlgorithm string

const (
	SSECustomerAlgorithmAES256 SSECustomerAlgorithm = "AES256"
)

type Owner struct {
	DisplayName *string
	ID          *string
}

type LocationConstraint string

type ObjectOwnership string

const (
	BucketOwnerPreferredOwnership ObjectOwnership = "BucketOwnerPreferred"
	ObjectWriterOwnership         ObjectOwnership = "ObjectWriter"
	BucketOwnerEnforcedOwnership  ObjectOwnership = "BucketOwnerEnforced"
)

type ChecksumMode string

const (
	ChecksumModeEnabled ChecksumMode = "ENABLED"
)

type ChecksumType string

const (
	ChecksumTypeComposite  ChecksumType = "COMPOSITE"
	ChecksumTypeFullObject ChecksumType = "FULL_OBJECT"
)

type SSEType string

const (
	SSE_AES256   SSEType = "AES256"
	SSE_KMS      SSEType = "aws:kms"
	SSE_KMS_DSSE SSEType = "aws:kms:dsse"
)

type StorageClass string

const (
	StorageClassStandard           StorageClass = "STANDARD"
	StorageClassReducedRedundancy  StorageClass = "REDUCED_REDUNDANCY"
	StorageClassStandardIA         StorageClass = "STANDARD_IA"
	StorageClassOnezoneIA          StorageClass = "ONEZONE_IA"
	StorageClassIntelligentTiering StorageClass = "INTELLIGENT_TIERING"
	StorageClassGlacier            StorageClass = "GLACIER"
	StorageClassDeepArchive        StorageClass = "DEEP_ARCHIVE"
	StorageClassOutposts           StorageClass = "OUTPOSTS"
	StorageClassGlacierIR          StorageClass = "GLACIER_IR"
	StorageClassSnow               StorageClass = "SNOW"
	StorageClassExpressOnezone     StorageClass = "EXPRESS_ONEZONE"
)

type ReplicationStatus string

const (
	ReplicationStatusComplete  = "COMPLETE"
	ReplicationStatusPending   = "PENDING"
	ReplicationStatusFailed    = "FAILED"
	ReplicationStatusReplica   = "REPLICA"
	ReplicationStatusCompleted = "COMPLETED"
)

type ObjectLockMode string

const (
	ObjectLockModeGovernance ObjectLockMode = "GOVERNANCE"
	ObjectLockModeCompliance ObjectLockMode = "COMPLIANCE"
)

type ObjectLockLegalHoldStatus string

const (
	ObjectLockLegalHoldStatusOn  ObjectLockLegalHoldStatus = "ON"
	ObjectLockLegalHoldStatusOff ObjectLockLegalHoldStatus = "OFF"
)

type ChecksumAlgorithm string

const (
	ChecksumAlgorithmCRC32     ChecksumAlgorithm = "CRC32"
	ChecksumAlgorithmCRC32C    ChecksumAlgorithm = "CRC32C"
	ChecksumAlgorithmCRC64NVME ChecksumAlgorithm = "CRC64NVME"
	ChecksumAlgorithmSHA1      ChecksumAlgorithm = "SHA1"
	ChecksumAlgorithmSHA256    ChecksumAlgorithm = "SHA256"
)
