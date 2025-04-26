package client

import (
	"errors"
)

type Client struct{}

func New() (*Client, error) {
	return nil, errors.ErrUnsupported
}

func (*Client) AbortMultipartUpload() error {
	return errors.ErrUnsupported
}

func (*Client) CompleteMultipartUpload() error {
	return errors.ErrUnsupported
}

func (*Client) CopyObject() error {
	return errors.ErrUnsupported
}

func (*Client) CreateBucket() error {
	return errors.ErrUnsupported
}

func (*Client) CreateBucketMetadataTableConfiguration() error {
	return errors.ErrUnsupported
}

func (*Client) CreateMultipartUpload() error {
	return errors.ErrUnsupported
}

func (*Client) CreateSession() error {
	return errors.ErrUnsupported
}

func (*Client) DeleteBucket() error {
	return errors.ErrUnsupported
}

func (*Client) DeleteBucketAnalyticsConfiguration() error {
	return errors.ErrUnsupported
}

func (*Client) DeleteBucketCors() error {
	return errors.ErrUnsupported
}

func (*Client) DeleteBucketEncryption() error {
	return errors.ErrUnsupported
}

func (*Client) DeleteBucketIntelligentTieringConfiguration() error {
	return errors.ErrUnsupported
}

func (*Client) DeleteBucketInventoryConfiguration() error {
	return errors.ErrUnsupported
}

func (*Client) DeleteBucketLifecycle() error {
	return errors.ErrUnsupported
}

func (*Client) DeleteBucketMetadataTableConfiguration() error {
	return errors.ErrUnsupported
}

func (*Client) DeleteBucketMetricsConfiguration() error {
	return errors.ErrUnsupported
}

func (*Client) DeleteBucketOwnershipControls() error {
	return errors.ErrUnsupported
}

func (*Client) DeleteBucketPolicy() error {
	return errors.ErrUnsupported
}

func (*Client) DeleteBucketReplication() error {
	return errors.ErrUnsupported
}

func (*Client) DeleteBucketTagging() error {
	return errors.ErrUnsupported
}

func (*Client) DeleteBucketWebsite() error {
	return errors.ErrUnsupported
}

func (*Client) DeleteObject() error {
	return errors.ErrUnsupported
}

func (*Client) DeleteObjects() error {
	return errors.ErrUnsupported
}

func (*Client) DeleteObjectTagging() error {
	return errors.ErrUnsupported
}

func (*Client) DeletePublicAccessBlock() error {
	return errors.ErrUnsupported
}

func (*Client) GetBucketAccelerateConfiguration() error {
	return errors.ErrUnsupported
}

func (*Client) GetBucketAcl() error {
	return errors.ErrUnsupported
}

func (*Client) GetBucketAnalyticsConfiguration() error {
	return errors.ErrUnsupported
}

func (*Client) GetBucketCors() error {
	return errors.ErrUnsupported
}

func (*Client) GetBucketEncryption() error {
	return errors.ErrUnsupported
}

func (*Client) GetBucketIntelligentTieringConfiguration() error {
	return errors.ErrUnsupported
}

func (*Client) GetBucketInventoryConfiguration() error {
	return errors.ErrUnsupported
}

func (*Client) GetBucketLifecycle() error {
	return errors.ErrUnsupported
}

func (*Client) GetBucketLifecycleConfiguration() error {
	return errors.ErrUnsupported
}

func (*Client) GetBucketLocation() error {
	return errors.ErrUnsupported
}

func (*Client) GetBucketLogging() error {
	return errors.ErrUnsupported
}

func (*Client) GetBucketMetadataTableConfiguration() error {
	return errors.ErrUnsupported
}

func (*Client) GetBucketMetricsConfiguration() error {
	return errors.ErrUnsupported
}

func (*Client) GetBucketNotification() error {
	return errors.ErrUnsupported
}

func (*Client) GetBucketNotificationConfiguration() error {
	return errors.ErrUnsupported
}

func (*Client) GetBucketOwnershipControls() error {
	return errors.ErrUnsupported
}

func (*Client) GetBucketPolicy() error {
	return errors.ErrUnsupported
}

func (*Client) GetBucketPolicyStatus() error {
	return errors.ErrUnsupported
}

func (*Client) GetBucketReplication() error {
	return errors.ErrUnsupported
}

func (*Client) GetBucketRequestPayment() error {
	return errors.ErrUnsupported
}

func (*Client) GetBucketTagging() error {
	return errors.ErrUnsupported
}

func (*Client) GetBucketVersioning() error {
	return errors.ErrUnsupported
}

func (*Client) GetBucketWebsite() error {
	return errors.ErrUnsupported
}

func (*Client) GetObject() error {
	return errors.ErrUnsupported
}

func (*Client) GetObjectAcl() error {
	return errors.ErrUnsupported
}

func (*Client) GetObjectAttributes() error {
	return errors.ErrUnsupported
}

func (*Client) GetObjectLegalHold() error {
	return errors.ErrUnsupported
}

func (*Client) GetObjectLockConfiguration() error {
	return errors.ErrUnsupported
}

func (*Client) GetObjectRetention() error {
	return errors.ErrUnsupported
}

func (*Client) GetObjectTagging() error {
	return errors.ErrUnsupported
}

func (*Client) GetObjectTorrent() error {
	return errors.ErrUnsupported
}

func (*Client) GetPublicAccessBlock() error {
	return errors.ErrUnsupported
}

func (*Client) HeadBucket() error {
	return errors.ErrUnsupported
}

func (*Client) HeadObject() error {
	return errors.ErrUnsupported
}

func (*Client) ListBucketAnalyticsConfigurations() error {
	return errors.ErrUnsupported
}

func (*Client) ListBucketIntelligentTieringConfigurations() error {
	return errors.ErrUnsupported
}

func (*Client) ListBucketInventoryConfigurations() error {
	return errors.ErrUnsupported
}

func (*Client) ListBucketMetricsConfigurations() error {
	return errors.ErrUnsupported
}

func (*Client) ListBuckets() error {
	return errors.ErrUnsupported
}

func (*Client) ListDirectoryBuckets() error {
	return errors.ErrUnsupported
}

func (*Client) ListMultipartUploads() error {
	return errors.ErrUnsupported
}

func (*Client) ListObjects() error {
	return errors.ErrUnsupported
}

func (*Client) ListObjectsV2() error {
	return errors.ErrUnsupported
}

func (*Client) ListObjectVersions() error {
	return errors.ErrUnsupported
}

func (*Client) ListParts() error {
	return errors.ErrUnsupported
}

func (*Client) PutBucketAccelerateConfiguration() error {
	return errors.ErrUnsupported
}

func (*Client) PutBucketAcl() error {
	return errors.ErrUnsupported
}

func (*Client) PutBucketAnalyticsConfiguration() error {
	return errors.ErrUnsupported
}

func (*Client) PutBucketCors() error {
	return errors.ErrUnsupported
}

func (*Client) PutBucketEncryption() error {
	return errors.ErrUnsupported
}

func (*Client) PutBucketIntelligentTieringConfiguration() error {
	return errors.ErrUnsupported
}

func (*Client) PutBucketInventoryConfiguration() error {
	return errors.ErrUnsupported
}

func (*Client) PutBucketLifecycle() error {
	return errors.ErrUnsupported
}

func (*Client) PutBucketLifecycleConfiguration() error {
	return errors.ErrUnsupported
}

func (*Client) PutBucketLogging() error {
	return errors.ErrUnsupported
}

func (*Client) PutBucketMetricsConfiguration() error {
	return errors.ErrUnsupported
}

func (*Client) PutBucketNotification() error {
	return errors.ErrUnsupported
}

func (*Client) PutBucketNotificationConfiguration() error {
	return errors.ErrUnsupported
}

func (*Client) PutBucketOwnershipControls() error {
	return errors.ErrUnsupported
}

func (*Client) PutBucketPolicy() error {
	return errors.ErrUnsupported
}

func (*Client) PutBucketReplication() error {
	return errors.ErrUnsupported
}

func (*Client) PutBucketRequestPayment() error {
	return errors.ErrUnsupported
}

func (*Client) PutBucketTagging() error {
	return errors.ErrUnsupported
}

func (*Client) PutBucketVersioning() error {
	return errors.ErrUnsupported
}

func (*Client) PutBucketWebsite() error {
	return errors.ErrUnsupported
}

func (*Client) PutObject() error {
	return errors.ErrUnsupported
}

func (*Client) PutObjectAcl() error {
	return errors.ErrUnsupported
}

func (*Client) PutObjectLegalHold() error {
	return errors.ErrUnsupported
}

func (*Client) PutObjectLockConfiguration() error {
	return errors.ErrUnsupported
}

func (*Client) PutObjectRetention() error {
	return errors.ErrUnsupported
}

func (*Client) PutObjectTagging() error {
	return errors.ErrUnsupported
}

func (*Client) PutPublicAccessBlock() error {
	return errors.ErrUnsupported
}

func (*Client) RestoreObject() error {
	return errors.ErrUnsupported
}

func (*Client) SelectObjectContent() error {
	return errors.ErrUnsupported
}

func (*Client) UploadPart() error {
	return errors.ErrUnsupported
}

func (*Client) UploadPartCopy() error {
	return errors.ErrUnsupported
}
