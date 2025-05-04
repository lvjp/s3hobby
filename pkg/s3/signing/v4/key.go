package v4

import (
	"slices"
	"strings"

	"github.com/lvjp/s3hobby/pkg/s3/signing"
	"github.com/lvjp/s3hobby/pkg/s3/signing/functions"
)

type SigningKey struct {
	v []byte
}

func NewSigningKey(credentials signing.Credentials, region string, signingTime signing.SigningTime) *SigningKey {
	dateKey := functions.HMAC_SHA256([]byte("AWS4"+credentials.SecretAccessKey), []byte(signingTime.ShortFormat()))
	dateRegionKey := functions.HMAC_SHA256(dateKey, []byte(region))
	dateRegionServiceKey := functions.HMAC_SHA256(dateRegionKey, []byte("s3"))

	return &SigningKey{
		v: functions.HMAC_SHA256(dateRegionServiceKey, []byte("aws4_request")),
	}
}

func (sk SigningKey) Sign(payload []byte) string {
	return functions.Hex(functions.HMAC_SHA256(sk.v, payload))
}

type StringToSignBuilder struct {
	signingTime string
	scope       string
}

func NewStringToSignBuilder(signingTime signing.SigningTime, region string) *StringToSignBuilder {
	return &StringToSignBuilder{
		signingTime: signingTime.LongFormat(),
		scope:       signingTime.ShortFormat() + "/" + region + "/s3/aws4_request",
	}
}

func (builder StringToSignBuilder) Scope() string {
	return builder.scope
}

func (builder StringToSignBuilder) BuildWith(method string, args ...string) string {
	return strings.Join(
		slices.Concat(
			[]string{
				method,
				builder.signingTime,
				builder.scope,
			},
			args,
		),
		"\n",
	)
}
