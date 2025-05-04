package v4

import (
	"bytes"
	"slices"
	"strconv"

	"github.com/lvjp/s3hobby/pkg/s3/api"
	"github.com/lvjp/s3hobby/pkg/s3/signing"
	"github.com/lvjp/s3hobby/pkg/s3/signing/functions"

	"github.com/valyala/fasthttp"
)

var emptyHash = functions.Hex(functions.SHA256Hash(nil))

type StreamedPayloadSigner struct {
	SignPayload       bool
	ForceEmptyTrailer bool
}

func (signer *StreamedPayloadSigner) Sign(args signing.SigningArgs) (canonicalRequest, stringToSign, signature string, err error) {
	trailer := signer.extractTrailer(args.Request)

	originalBody := args.Request.Body()
	encodedLength := GetStreamEncodedContentLength(signer.SignPayload, originalBody, trailer)

	args.Request.Header.Set(api.HeaderContentEncoding, signer.computeContentEncoding(args.Request))
	args.Request.Header.Set(api.HeaderXAmzContentSHA256, signer.computePayloadSHA256Hash(trailer != nil))
	args.Request.Header.SetContentLength(encodedLength)
	args.Request.Header.Set(api.HeaderXAmzDecodedContentLength, strconv.Itoa(len(originalBody)))
	args.Request.Header.Set(api.HeaderXAmzDate, args.SigningTime.LongFormat())

	var authorizationHeader string
	canonicalRequest, stringToSign, signature, authorizationHeader, err = getHeaderSignature(args)
	if err != nil {
		return "", "", "", err
	}

	args.Request.Header.Set(api.HeaderAuthorization, authorizationHeader)

	var payloadSigner *StreamPayloadSigner

	if signer.SignPayload {
		payloadSigner = NewStreamPayloadSigner(
			NewSigningKey(args.Credentials, args.Region, args.SigningTime),
			signature,
			NewStringToSignBuilder(args.SigningTime, args.Region),
		)
	}

	args.Request.SetBodyRaw(signer.transform(payloadSigner, encodedLength, originalBody, trailer))

	return
}

func (*StreamedPayloadSigner) computeContentEncoding(req *fasthttp.Request) string {
	contentEncoding := "aws-chunked"
	if actualContentEncoding := req.Header.Peek(api.HeaderContentEncoding); len(actualContentEncoding) > 0 {
		contentEncoding += "," + string(actualContentEncoding)
	}

	return contentEncoding
}

func (signer *StreamedPayloadSigner) computePayloadSHA256Hash(haveTrailer bool) string {
	switch {
	case !signer.SignPayload && !haveTrailer:
		// Will send an empty trailer
		return "STREAMING-UNSIGNED-PAYLOAD-TRAILER"
	case !signer.SignPayload && haveTrailer:
		return "STREAMING-UNSIGNED-PAYLOAD-TRAILER"
	case signer.SignPayload && !haveTrailer:
		return "STREAMING-AWS4-HMAC-SHA256-PAYLOAD"
	case signer.SignPayload && haveTrailer:
		return "STREAMING-AWS4-HMAC-SHA256-PAYLOAD-TRAILER"
	default:
		panic("unreachable")
	}
}

func (signer *StreamedPayloadSigner) extractTrailer(req *fasthttp.Request) *TrailerBody {
	var trailer *TrailerBody

	if trailerName := req.Header.Peek(api.HeaderXAmzTrailer); len(trailerName) > 0 {
		trailer = &TrailerBody{
			Name:  string(trailerName),
			Value: string(req.Header.PeekBytes(trailerName)),
		}
		req.Header.DelBytes(trailerName)
	} else if signer.ForceEmptyTrailer {
		trailer = &TrailerBody{}
	}

	return trailer
}

func (signer *StreamedPayloadSigner) transform(payloadSigner *StreamPayloadSigner, encodedLen int, originalBody []byte, trailer *TrailerBody) []byte {
	var buf bytes.Buffer
	buf.Grow(encodedLen)

	for chunk := range slices.Chunk(originalBody, chunkDataSize) {
		signer.writeChunk(payloadSigner, &buf, chunk)
	}

	signer.writeChunk(payloadSigner, &buf, nil)

	if trailer != nil {
		buf.Write(trailer.Bytes())
		if payloadSigner != nil {
			_, signature := payloadSigner.TrailerSignature([]byte(trailer.StringtoSign()))
			buf.WriteString(api.HeaderXAmzTrailerSignature)
			buf.WriteString(trailerSeparator)
			buf.WriteString(signature)
			buf.Write(crlf)
		}
	}

	buf.Write(crlf)

	return buf.Bytes()
}

func (*StreamedPayloadSigner) writeChunk(payloadSigner *StreamPayloadSigner, body *bytes.Buffer, chunk []byte) {
	if payloadSigner == nil {
		body.WriteString(strconv.FormatInt(int64(len(chunk)), 16))
		body.WriteString("\r\n")
		if len(chunk) > 0 {
			body.Write(chunk)
			body.WriteString("\r\n")
		}
		return
	}

	_, signature := payloadSigner.ChunkSignature(chunk)

	body.WriteString(strconv.FormatInt(int64(len(chunk)), 16))
	body.WriteString(";chunk-signature=")
	body.WriteString(signature)
	body.WriteString("\r\n")
	if len(chunk) > 0 {
		body.Write(chunk)
		body.WriteString("\r\n")
	}
}

type StreamPayloadSigner struct {
	signingKey        *SigningKey
	previousSignature string
	stringToSign      *StringToSignBuilder
}

func NewStreamPayloadSigner(signingKey *SigningKey, seedSignature string, builder *StringToSignBuilder) *StreamPayloadSigner {
	return &StreamPayloadSigner{
		signingKey:        signingKey,
		previousSignature: seedSignature,
		stringToSign:      builder,
	}
}

func (s *StreamPayloadSigner) ChunkSignature(payload []byte) (stringToSign, signature string) {
	stringToSign = s.stringToSign.BuildWith(
		"AWS4-HMAC-SHA256-PAYLOAD",
		s.previousSignature,
		emptyHash,
		functions.Hex(functions.SHA256Hash(payload)),
	)

	signature = s.signingKey.Sign([]byte(stringToSign))
	s.previousSignature = signature

	return
}

func (s *StreamPayloadSigner) TrailerSignature(payload []byte) (stringToSign, signature string) {
	stringToSign = s.stringToSign.BuildWith(
		"AWS4-HMAC-SHA256-TRAILER",
		s.previousSignature,
		functions.Hex(functions.SHA256Hash(payload)),
	)

	signature = s.signingKey.Sign([]byte(stringToSign))
	s.previousSignature = signature

	return
}

const chunkDataSize = 64 * 1024
const trailerSeparator = ":"

var crlf = []byte{'\r', '\n'}

type TrailerBody struct {
	Name  string
	Value string
}

func (t *TrailerBody) Len() int {
	if len(t.Name) == 0 {
		return 0
	}

	return len(t.Name) + len(trailerSeparator) + len(t.Value) + len(crlf)
}

func (t *TrailerBody) StringtoSign() string {
	if len(t.Name) == 0 {
		return "\n"
	}

	return t.Name + trailerSeparator + t.Value + "\n"
}

func (t *TrailerBody) Bytes() []byte {
	var buf bytes.Buffer
	buf.Grow(t.Len())

	var signBody string

	if len(t.Name) > 0 {
		signBody = t.Name + trailerSeparator + t.Value
		buf.WriteString(signBody)
		buf.Write(crlf)
	}

	return buf.Bytes()
}

func GetStreamEncodedContentLength(isSigned bool, originalBody []byte, trailer *TrailerBody) int {
	const signatureValueLen = 64
	const signatureSize = len(";chunk-signature=") + signatureValueLen

	decodedContentLength := len(originalBody)

	bodyLen := decodedContentLength
	if nbChunk := decodedContentLength / chunkDataSize; nbChunk > 0 {
		chunkSize := len(strconv.FormatInt(chunkDataSize, 16)) + len(crlf)
		if isSigned {
			chunkSize += signatureSize
		}
		chunkSize += len(crlf)

		bodyLen += nbChunk * chunkSize
	}

	if remaining := decodedContentLength % chunkDataSize; remaining > 0 {
		bodyLen += len(strconv.FormatInt(int64(remaining), 16)) + len(crlf)
		if isSigned {
			bodyLen += signatureSize
		}
		bodyLen += len(crlf)
	}

	bodyLen += len("0") + len(crlf)
	if isSigned {
		bodyLen += signatureSize
	}

	if trailer != nil {
		bodyLen += trailer.Len()
		if isSigned {
			bodyLen += len(api.HeaderXAmzTrailerSignature) + len(trailerSeparator) + signatureValueLen + len(crlf)
		}
	}

	bodyLen += len(crlf)

	return bodyLen
}
