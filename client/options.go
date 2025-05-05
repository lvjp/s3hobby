package client

import (
	"github.com/lvjp/s3hobby/pkg/s3/signing"

	"github.com/go-playground/validator/v10"
)

const DefaultUserAgent = "s3hobby-client"

type Options struct {
	// UserAgent specifies how to populate the User-Agent api.Header.
	// [DefaultUserAgent] is used when nil.
	// Nothing is sent when the pointed value is empty.
	// Otherwise, the pointed value is sent.
	UserAgent *string

	UsePathStyle bool
	EndpointHost string `validate:"hostname|hostname_port"`
	UseSSL       bool

	// EndpointResolver default to [DefaultEndpointResolver].
	EndpointResolver EndpointResolver `validate:"required"`

	SiginingRegion string `validate:"required"`

	Signer signing.Signer `validate:"required"`

	Credentials *signing.Credentials

	// HTTPClient default to [DefaultHTTPClient].
	HTTPClient HTTPClient `validate:"required"`
}

// With return a new instance of [Options] with applied transformations.
func (opts *Options) With(optFns ...func(*Options)) *Options {
	ret := *opts

	for _, fn := range optFns {
		fn(&ret)
	}

	return &ret
}

func (opts *Options) Validate() error {
	return validator.New(validator.WithRequiredStructEnabled()).Struct(opts)
}

func (opts *Options) setDefaults() {
	if opts.UserAgent == nil {
		userAgent := DefaultUserAgent
		opts.UserAgent = &userAgent
	}

	if opts.EndpointResolver == nil {
		opts.EndpointResolver = &DefaultEndpointResolver{}
	}

	if opts.HTTPClient == nil {
		opts.HTTPClient = DefaultHTTPClient
	}
}
