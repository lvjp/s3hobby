package anonymous

import "github.com/lvjp/s3hobby/pkg/s3/signing"

func NewSigner() signing.Signer {
	return &anonymousImpl{}
}

type anonymousImpl struct {
}

func (*anonymousImpl) Sign(signing.SigningArgs) (canonicalRequest, stringToSign, signature string, err error) {
	return "", "", "", nil
}
