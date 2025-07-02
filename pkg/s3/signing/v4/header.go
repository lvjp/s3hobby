package v4

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/lvjp/s3hobby/pkg/s3/api"
	"github.com/lvjp/s3hobby/pkg/s3/signing"
	"github.com/lvjp/s3hobby/pkg/s3/signing/functions"
)

type PlainPayloadSigner struct {
	SignPayload bool
}

func (signer *PlainPayloadSigner) Sign(args signing.SigningArgs) (canonicalRequest, stringToSign, signature string, err error) {
	hashPayload := "UNSIGNED-PAYLOAD"
	if signer.SignPayload {
		hashPayload = functions.Hex(functions.SHA256Hash(args.Request.Body()))
	}
	args.Request.Header.Set(api.HeaderXAmzContentSHA256, hashPayload)
	args.Request.Header.Set(api.HeaderXAmzDate, args.SigningTime.LongFormat())

	var authorizationHeader string
	canonicalRequest, stringToSign, signature, authorizationHeader, err = getHeaderSignature(args)
	if err != nil {
		return "", "", "", err
	}

	args.Request.Header.Set(api.HeaderAuthorization, authorizationHeader)

	return
}

func getHeaderSignature(args signing.SigningArgs) (canonicalRequest, stringToSign, signature, authorizationHeader string, err error) {
	ctx := newHeaderSigningCtx(args)

	if actual, expected := ctx.Request.Header.Peek(api.HeaderXAmzDate), ctx.SigningTime.LongFormat(); string(actual) != expected {
		return "", "", "", "", fmt.Errorf("HeaderSigner: %q header mismatch: expected %q, got %q", ctx.SigningTime.LongFormat(), expected, actual)
	}

	payloadHash := string(ctx.Request.Header.Peek(api.HeaderXAmzContentSHA256))
	if len(payloadHash) == 0 {
		return "", "", "", "", fmt.Errorf("HeaderSigner: %q header not found", api.HeaderXAmzContentSHA256)
	}

	canonicalHeaders, signedHeaders := ctx.computeHeaders()
	canonicalRequest = ctx.computeCanonicalRequest(canonicalHeaders, signedHeaders, payloadHash)

	method := "AWS4-HMAC-SHA256"
	stringToSignBuilder := NewStringToSignBuilder(ctx.SigningTime, ctx.Region)
	stringToSign = stringToSignBuilder.BuildWith(method, functions.Hex(functions.SHA256Hash([]byte(canonicalRequest))))
	signature = ctx.signingKey.Sign([]byte(stringToSign))

	authorizationHeader = fmt.Sprintf(
		"%s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		method,
		ctx.Credentials.AccessKeyID,
		stringToSignBuilder.Scope(),
		signedHeaders,
		signature,
	)

	return canonicalRequest, stringToSign, signature, authorizationHeader, nil
}

type headerSigningCtx struct {
	signing.SigningArgs

	signingKey *SigningKey
}

func newHeaderSigningCtx(args signing.SigningArgs) *headerSigningCtx {
	return &headerSigningCtx{
		SigningArgs: args,
		signingKey: NewSigningKey(
			args.Credentials,
			args.Region,
			args.SigningTime,
		),
	}
}

func (ctx *headerSigningCtx) computeHeaders() (canonicalHeaders, signedHeaders string) {
	normalized := make(map[string]string, ctx.Request.Header.Len())

	for key, value := range ctx.Request.Header.All() {
		normalizedKey := functions.LowerCase(string(key))
		normalizedValue := functions.Trim(string(value))

		normalized[normalizedKey] = normalizedValue
	}

	// Ensure host header presence
	if _, exists := normalized["host"]; !exists {
		normalized["host"] = string(ctx.Request.Host())
	}

	sortedHeaders := slices.Sorted(maps.Keys(normalized))
	for _, key := range sortedHeaders {
		canonicalHeaders += key + ":" + normalized[key] + "\n"
	}

	signedHeaders = strings.Join(sortedHeaders, ";")

	return
}

func (ctx *headerSigningCtx) computeCanonicalRequest(canonicalHeaders, signedHeaders, payloadHash string) string {
	path := string(ctx.Request.URI().PathOriginal())
	if path == "" {
		path = "/"
	}

	return strings.Join(
		[]string{
			string(ctx.Request.Header.Method()),
			functions.URIEncode(path, true),
			ctx.getCanonicalQueryString(),
			canonicalHeaders,
			signedHeaders,
			payloadHash,
		},
		"\n",
	)
}

func (ctx *headerSigningCtx) getCanonicalQueryString() string {
	args := ctx.Request.URI().QueryArgs()
	encoded := make(map[string]string, args.Len())

	for key, value := range args.All() {
		encoded[functions.URIEncode(string(key), false)] = functions.URIEncode(string(value), false)
	}

	var ret strings.Builder

	for i, key := range slices.Sorted(maps.Keys(encoded)) {
		if i > 0 {
			ret.WriteRune('&')
		}

		ret.WriteString(key)
		ret.WriteRune('=')
		ret.WriteString(encoded[key])
	}

	return ret.String()
}
