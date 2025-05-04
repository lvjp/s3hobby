package v4

import (
	"bytes"

	"github.com/lvjp/s3hobby/pkg/s3/api"
	"github.com/lvjp/s3hobby/pkg/s3/signing"
)

type dynamicSigner struct{}

func NewDynamicSigner() signing.Signer {
	return &dynamicSigner{}
}

func (*dynamicSigner) Sign(args signing.SigningArgs) (canonicalRequest, stringToSign, signature string, err error) {
	needStream := len(args.Request.Header.Peek(api.HeaderXAmzTrailer)) > 0
	signPayload := bytes.Equal(args.Request.URI().Scheme(), []byte("http"))

	var signer signing.Signer = &PlainPayloadSigner{SignPayload: signPayload}
	if needStream {
		signer = &StreamedPayloadSigner{SignPayload: signPayload}
	}

	return signer.Sign(args)
}

func NewSignerWith(signPayload, streamPayload, forceEmptyTrailer bool) signing.Signer {
	if streamPayload || forceEmptyTrailer {
		return &StreamedPayloadSigner{
			SignPayload:       signPayload,
			ForceEmptyTrailer: forceEmptyTrailer,
		}
	}

	return &PlainPayloadSigner{
		SignPayload: signPayload,
	}
}
