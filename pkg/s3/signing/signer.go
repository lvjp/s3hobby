package signing

import (
	"github.com/valyala/fasthttp"
)

type Credentials struct {
	AccessKeyID     string
	SecretAccessKey string
}

type SigningArgs struct {
	Request *fasthttp.Request

	Credentials Credentials
	Region      string
	SigningTime SigningTime
}

type Signer interface {
	Sign(args SigningArgs) (canonicalRequest, stringToSign, signature string, err error)
}
