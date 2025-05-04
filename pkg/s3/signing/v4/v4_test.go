package v4

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/lvjp/s3hobby/pkg/s3/api"
	"github.com/lvjp/s3hobby/pkg/s3/signing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

func TestSigners(t *testing.T) {
	const URI = "https://examplebucket.s3.amazonaws.com/photos/photo1.jpg"

	type input struct {
		builder func(req *fasthttp.Request)
		signer  signing.Signer
	}

	type expected struct {
		canonicalRequest []string
		stringToSign     []string
		signature        string

		headers []string
		body    []string
	}

	testCases := []struct {
		name string

		input    input
		expected expected
	}{
		{
			name: "UNSIGNED-PAYLOAD/no payload",
			input: input{
				builder: func(req *fasthttp.Request) {
					req.Header.SetMethod("PUT")
					req.SetRequestURI(URI)
				},
				signer: &PlainPayloadSigner{},
			},
			expected: expected{
				canonicalRequest: []string{
					"PUT",
					"/photos/photo1.jpg",
					"",
					"host:examplebucket.s3.amazonaws.com",
					"x-amz-content-sha256:UNSIGNED-PAYLOAD",
					"x-amz-date:19840805T135000Z",
					"",
					"host;x-amz-content-sha256;x-amz-date",
					"UNSIGNED-PAYLOAD",
				},
				stringToSign: []string{
					"AWS4-HMAC-SHA256",
					"19840805T135000Z",
					"19840805/us-east-1/s3/aws4_request",
					"86fe1777c26bd70ba769eaa89a46aa69d67ea398ac096549a1cf274562940577",
				},
				signature: "20c9433818855c36457c61c41f3ca58ad5c368ed0b2c62a46e60c46be6f75a61",
				headers: []string{
					"PUT https://examplebucket.s3.amazonaws.com/photos/photo1.jpg HTTP/1.1",
					"X-Amz-Content-Sha256: UNSIGNED-PAYLOAD",
					"X-Amz-Date: 19840805T135000Z",
					"Authorization: AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/19840805/us-east-1/s3/aws4_request, SignedHeaders=host;x-amz-content-sha256;x-amz-date, Signature=20c9433818855c36457c61c41f3ca58ad5c368ed0b2c62a46e60c46be6f75a61",
				},
			},
		},
		{
			name: "UNSIGNED-PAYLOAD/with payload",
			input: input{
				builder: func(req *fasthttp.Request) {
					req.Header.SetMethod("PUT")
					req.SetRequestURI(URI)
					req.Header.Set(api.HeaderXAmzChecksumCrc64nvme, "ntuPBsmdl18=")
					req.SetBody([]byte("Welcome to S3."))
				},
				signer: &PlainPayloadSigner{},
			},
			expected: expected{
				canonicalRequest: []string{
					"PUT",
					"/photos/photo1.jpg",
					"",
					"host:examplebucket.s3.amazonaws.com",
					"x-amz-checksum-crc64nvme:ntuPBsmdl18=",
					"x-amz-content-sha256:UNSIGNED-PAYLOAD",
					"x-amz-date:19840805T135000Z",
					"",
					"host;x-amz-checksum-crc64nvme;x-amz-content-sha256;x-amz-date",
					"UNSIGNED-PAYLOAD",
				},
				stringToSign: []string{
					"AWS4-HMAC-SHA256",
					"19840805T135000Z",
					"19840805/us-east-1/s3/aws4_request",
					"b0b803fc2e4edcc2e2af79a15a3cb428b11e8368340607735bcd44d483e65100",
				},
				signature: "deddf092b5828abac288677ffc7a911baf49c743b7888577490f36a3784ab2f0",
				headers: []string{
					"PUT https://examplebucket.s3.amazonaws.com/photos/photo1.jpg HTTP/1.1",
					"X-Amz-Checksum-Crc64nvme: ntuPBsmdl18=",
					"X-Amz-Content-Sha256: UNSIGNED-PAYLOAD",
					"X-Amz-Date: 19840805T135000Z",
					"Authorization: AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/19840805/us-east-1/s3/aws4_request, SignedHeaders=host;x-amz-checksum-crc64nvme;x-amz-content-sha256;x-amz-date, Signature=deddf092b5828abac288677ffc7a911baf49c743b7888577490f36a3784ab2f0",
				},
				body: []string{
					"Welcome to S3.",
				},
			},
		},
		{
			name: "SIGNED-PAYLOAD/no payload",
			input: input{
				builder: func(req *fasthttp.Request) {
					req.Header.SetMethod("PUT")
					req.SetRequestURI(URI)
				},
				signer: &PlainPayloadSigner{SignPayload: true},
			},
			expected: expected{
				canonicalRequest: []string{
					"PUT",
					"/photos/photo1.jpg",
					"",
					"host:examplebucket.s3.amazonaws.com",
					"x-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
					"x-amz-date:19840805T135000Z",
					"",
					"host;x-amz-content-sha256;x-amz-date",
					"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				},
				stringToSign: []string{
					"AWS4-HMAC-SHA256",
					"19840805T135000Z",
					"19840805/us-east-1/s3/aws4_request",
					"155b8102f0017eb35b2121da423b904735dc5e1098e7092a7f3c7eb165ac6583",
				},
				signature: "2d0131f8c82108681757d7da2f492040b803b0c51bd24d631012ab79e3932ddb",
				headers: []string{
					"PUT https://examplebucket.s3.amazonaws.com/photos/photo1.jpg HTTP/1.1",
					"X-Amz-Content-Sha256: e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
					"X-Amz-Date: 19840805T135000Z",
					"Authorization: AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/19840805/us-east-1/s3/aws4_request, SignedHeaders=host;x-amz-content-sha256;x-amz-date, Signature=2d0131f8c82108681757d7da2f492040b803b0c51bd24d631012ab79e3932ddb",
				},
			},
		},
		{
			name: "SIGNED-PAYLOAD/with payload",
			input: input{
				builder: func(req *fasthttp.Request) {
					req.Header.SetMethod("PUT")
					req.SetRequestURI(URI + "?x-id=PutObject")
					req.Header.Set(api.HeaderXAmzChecksumCrc64nvme, "ntuPBsmdl18=")
					req.SetBody([]byte("Welcome to S3."))
				},
				signer: &PlainPayloadSigner{SignPayload: true},
			},
			expected: expected{
				canonicalRequest: []string{
					"PUT",
					"/photos/photo1.jpg",
					"x-id=PutObject",
					"host:examplebucket.s3.amazonaws.com",
					"x-amz-checksum-crc64nvme:ntuPBsmdl18=",
					"x-amz-content-sha256:f3893d4cc3e907c99afd2b35ae83e391b914b78c98097d9b5f7c89d4800fbaa9",
					"x-amz-date:19840805T135000Z",
					"",
					"host;x-amz-checksum-crc64nvme;x-amz-content-sha256;x-amz-date",
					"f3893d4cc3e907c99afd2b35ae83e391b914b78c98097d9b5f7c89d4800fbaa9",
				},
				stringToSign: []string{
					"AWS4-HMAC-SHA256",
					"19840805T135000Z",
					"19840805/us-east-1/s3/aws4_request",
					"95fb6985ad7211670de41749892e7e1caf8326f9c2c9c95ff7cd8e8c04c9a9cc",
				},
				signature: "d788eaeeb9ee7b1fe3ace0cd5bbf1d122b498a54f3bbdb6273c97a31cb1171f8",
				headers: []string{
					"PUT https://examplebucket.s3.amazonaws.com/photos/photo1.jpg?x-id=PutObject HTTP/1.1",
					"X-Amz-Checksum-Crc64nvme: ntuPBsmdl18=",
					"X-Amz-Content-Sha256: f3893d4cc3e907c99afd2b35ae83e391b914b78c98097d9b5f7c89d4800fbaa9",
					"X-Amz-Date: 19840805T135000Z",
					"Authorization: AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/19840805/us-east-1/s3/aws4_request, SignedHeaders=host;x-amz-checksum-crc64nvme;x-amz-content-sha256;x-amz-date, Signature=d788eaeeb9ee7b1fe3ace0cd5bbf1d122b498a54f3bbdb6273c97a31cb1171f8",
				},
				body: []string{
					"Welcome to S3.",
				},
			},
		},
		{
			name: "STREAMING-UNSIGNED-PAYLOAD-TRAILER/no payload/empty trailer",
			input: input{
				builder: func(req *fasthttp.Request) {
					req.Header.SetMethod("PUT")
					req.SetRequestURI(URI + "?x-id=PutObject")
				},
				signer: &StreamedPayloadSigner{},
			},
			expected: expected{
				canonicalRequest: []string{
					"PUT",
					"/photos/photo1.jpg",
					"x-id=PutObject",
					"content-encoding:aws-chunked",
					"content-length:5",
					"host:examplebucket.s3.amazonaws.com",
					"x-amz-content-sha256:STREAMING-UNSIGNED-PAYLOAD-TRAILER",
					"x-amz-date:19840805T135000Z",
					"x-amz-decoded-content-length:0",
					"",
					"content-encoding;content-length;host;x-amz-content-sha256;x-amz-date;x-amz-decoded-content-length",
					"STREAMING-UNSIGNED-PAYLOAD-TRAILER",
				},
				stringToSign: []string{
					"AWS4-HMAC-SHA256",
					"19840805T135000Z",
					"19840805/us-east-1/s3/aws4_request",
					"7ea79cdbd075b8fc441529b39c7747fd9966a5a9f52138b57068a78f994faddd",
				},
				signature: "404739ee845b163c5fdb8c39209c0d806eba9248b7ed876f608471addc3b41a5",
				headers: []string{
					"PUT https://examplebucket.s3.amazonaws.com/photos/photo1.jpg?x-id=PutObject HTTP/1.1",
					"Content-Length: 5",
					"Content-Encoding: aws-chunked",
					"X-Amz-Content-Sha256: STREAMING-UNSIGNED-PAYLOAD-TRAILER",
					"X-Amz-Decoded-Content-Length: 0",
					"X-Amz-Date: 19840805T135000Z",
					"Authorization: AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/19840805/us-east-1/s3/aws4_request, SignedHeaders=content-encoding;content-length;host;x-amz-content-sha256;x-amz-date;x-amz-decoded-content-length, Signature=404739ee845b163c5fdb8c39209c0d806eba9248b7ed876f608471addc3b41a5",
				},
				body: []string{
					"0\r\n",
					"\r\n",
				},
			},
		},
		{
			name: "STREAMING-UNSIGNED-PAYLOAD-TRAILER/no payload/with trailer",
			input: input{
				builder: func(req *fasthttp.Request) {
					req.Header.SetMethod("PUT")
					req.SetRequestURI(URI + "?x-id=PutObject")
					req.Header.Set(api.HeaderXAmzTrailer, api.HeaderXAmzChecksumCrc64nvme)
					req.Header.Set(api.HeaderXAmzChecksumCrc64nvme, "AAAAAAAAAAA=")
				},
				signer: &StreamedPayloadSigner{},
			},
			expected: expected{
				canonicalRequest: []string{
					"PUT",
					"/photos/photo1.jpg",
					"x-id=PutObject",
					"content-encoding:aws-chunked",
					"content-length:44",
					"host:examplebucket.s3.amazonaws.com",
					"x-amz-content-sha256:STREAMING-UNSIGNED-PAYLOAD-TRAILER",
					"x-amz-date:19840805T135000Z",
					"x-amz-decoded-content-length:0",
					"x-amz-trailer:x-amz-checksum-crc64nvme",
					"",
					"content-encoding;content-length;host;x-amz-content-sha256;x-amz-date;x-amz-decoded-content-length;x-amz-trailer",
					"STREAMING-UNSIGNED-PAYLOAD-TRAILER",
				},
				stringToSign: []string{
					"AWS4-HMAC-SHA256",
					"19840805T135000Z",
					"19840805/us-east-1/s3/aws4_request",
					"47ea6fb330bbe15b871322f60d507f0f5530e5406f1c2d25606223bb2ccf1821",
				},
				signature: "7b520fe22d801e430f74e2c8c1c712b16d86dc69d4f268e5ec0dae093407b5f3",
				headers: []string{
					"PUT https://examplebucket.s3.amazonaws.com/photos/photo1.jpg?x-id=PutObject HTTP/1.1",
					"Content-Length: 44",
					"X-Amz-Trailer: x-amz-checksum-crc64nvme",
					"Content-Encoding: aws-chunked",
					"X-Amz-Content-Sha256: STREAMING-UNSIGNED-PAYLOAD-TRAILER",
					"X-Amz-Decoded-Content-Length: 0",
					"X-Amz-Date: 19840805T135000Z",
					"Authorization: AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/19840805/us-east-1/s3/aws4_request, SignedHeaders=content-encoding;content-length;host;x-amz-content-sha256;x-amz-date;x-amz-decoded-content-length;x-amz-trailer, Signature=7b520fe22d801e430f74e2c8c1c712b16d86dc69d4f268e5ec0dae093407b5f3",
				},
				body: []string{
					"0\r\n",
					"x-amz-checksum-crc64nvme:AAAAAAAAAAA=\r\n",
					"\r\n",
				},
			},
		},
		{
			name: "STREAMING-UNSIGNED-PAYLOAD-TRAILER/with payload/empty trailer",
			input: input{
				builder: func(req *fasthttp.Request) {
					req.Header.SetMethod("PUT")
					req.SetRequestURI(URI + "?x-id=PutObject")
					req.SetBody(bytes.Repeat([]byte("a"), 64*1024+1024))
				},
				signer: &StreamedPayloadSigner{ForceEmptyTrailer: true},
			},
			expected: expected{
				canonicalRequest: []string{
					"PUT",
					"/photos/photo1.jpg",
					"x-id=PutObject",
					"content-encoding:aws-chunked",
					"content-length:66581",
					"host:examplebucket.s3.amazonaws.com",
					"x-amz-content-sha256:STREAMING-UNSIGNED-PAYLOAD-TRAILER",
					"x-amz-date:19840805T135000Z",
					"x-amz-decoded-content-length:66560",
					"",
					"content-encoding;content-length;host;x-amz-content-sha256;x-amz-date;x-amz-decoded-content-length",
					"STREAMING-UNSIGNED-PAYLOAD-TRAILER",
				},
				stringToSign: []string{
					"AWS4-HMAC-SHA256",
					"19840805T135000Z",
					"19840805/us-east-1/s3/aws4_request",
					"00fb094d2d08f19aa8a6596bc2ad5eebdbfda8dce0fe2b5b694786ff18de62c8",
				},
				signature: "3b085d94c65e07c780503a9224f5a40474eebd759f0c1bc7a516d7f270a9c8b8",
				headers: []string{
					"PUT https://examplebucket.s3.amazonaws.com/photos/photo1.jpg?x-id=PutObject HTTP/1.1",
					"Content-Length: 66581",
					"Content-Encoding: aws-chunked",
					"X-Amz-Content-Sha256: STREAMING-UNSIGNED-PAYLOAD-TRAILER",
					"X-Amz-Decoded-Content-Length: 66560",
					"X-Amz-Date: 19840805T135000Z",
					"Authorization: AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/19840805/us-east-1/s3/aws4_request, SignedHeaders=content-encoding;content-length;host;x-amz-content-sha256;x-amz-date;x-amz-decoded-content-length, Signature=3b085d94c65e07c780503a9224f5a40474eebd759f0c1bc7a516d7f270a9c8b8",
				},
				body: []string{
					"10000\r\n",
					strings.Repeat("a", 64*1024) + "\r\n",
					"400\r\n",
					strings.Repeat("a", 1024) + "\r\n",
					"0\r\n",
					"\r\n",
				},
			},
		},
		{
			name: "STREAMING-UNSIGNED-PAYLOAD-TRAILER/with payload/with trailer",
			input: input{
				builder: func(req *fasthttp.Request) {
					req.Header.SetMethod("PUT")
					req.SetRequestURI(URI + "?x-id=PutObject")
					req.SetBody(bytes.Repeat([]byte("a"), 64*1024+1024))
					req.Header.Set(api.HeaderXAmzTrailer, api.HeaderXAmzChecksumCrc32c)
					req.Header.Set(api.HeaderXAmzChecksumCrc32c, "sOO8/Q==")
				},
				signer: &StreamedPayloadSigner{},
			},
			expected: expected{
				canonicalRequest: []string{
					"PUT",
					"/photos/photo1.jpg",
					"x-id=PutObject",
					"content-encoding:aws-chunked",
					"content-length:66613",
					"host:examplebucket.s3.amazonaws.com",
					"x-amz-content-sha256:STREAMING-UNSIGNED-PAYLOAD-TRAILER",
					"x-amz-date:19840805T135000Z",
					"x-amz-decoded-content-length:66560",
					"x-amz-trailer:x-amz-checksum-crc32c",
					"",
					"content-encoding;content-length;host;x-amz-content-sha256;x-amz-date;x-amz-decoded-content-length;x-amz-trailer",
					"STREAMING-UNSIGNED-PAYLOAD-TRAILER",
				},
				stringToSign: []string{
					"AWS4-HMAC-SHA256",
					"19840805T135000Z",
					"19840805/us-east-1/s3/aws4_request",
					"2cb918f54cafda3736ac611b2ed772373e1ed6396f90868e1f61d88e3d7b3439",
				},
				signature: "c153df61a82ba69aa920fa50bc32091d20507e491b7883db976ecee4e82c38cf",
				headers: []string{
					"PUT https://examplebucket.s3.amazonaws.com/photos/photo1.jpg?x-id=PutObject HTTP/1.1",
					"Content-Length: 66613",
					"X-Amz-Trailer: x-amz-checksum-crc32c",
					"Content-Encoding: aws-chunked",
					"X-Amz-Content-Sha256: STREAMING-UNSIGNED-PAYLOAD-TRAILER",
					"X-Amz-Decoded-Content-Length: 66560",
					"X-Amz-Date: 19840805T135000Z",
					"Authorization: AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/19840805/us-east-1/s3/aws4_request, SignedHeaders=content-encoding;content-length;host;x-amz-content-sha256;x-amz-date;x-amz-decoded-content-length;x-amz-trailer, Signature=c153df61a82ba69aa920fa50bc32091d20507e491b7883db976ecee4e82c38cf",
				},
				body: []string{
					"10000\r\n",
					strings.Repeat("a", 64*1024) + "\r\n",
					"400\r\n",
					strings.Repeat("a", 1024) + "\r\n",
					"0\r\n",
					"x-amz-checksum-crc32c:sOO8/Q==\r\n",
					"\r\n",
				},
			},
		},
		{
			name: "STREAMING-AWS4-HMAC-SHA256-PAYLOAD/no payload",
			input: input{
				builder: func(req *fasthttp.Request) {
					req.Header.SetMethod("PUT")
					req.SetRequestURI(URI)
				},
				signer: &StreamedPayloadSigner{SignPayload: true},
			},
			expected: expected{
				canonicalRequest: []string{
					"PUT",
					"/photos/photo1.jpg",
					"",
					"content-encoding:aws-chunked",
					"content-length:86",
					"host:examplebucket.s3.amazonaws.com",
					"x-amz-content-sha256:STREAMING-AWS4-HMAC-SHA256-PAYLOAD",
					"x-amz-date:19840805T135000Z",
					"x-amz-decoded-content-length:0",
					"",
					"content-encoding;content-length;host;x-amz-content-sha256;x-amz-date;x-amz-decoded-content-length",
					"STREAMING-AWS4-HMAC-SHA256-PAYLOAD",
				},
				stringToSign: []string{
					"AWS4-HMAC-SHA256",
					"19840805T135000Z",
					"19840805/us-east-1/s3/aws4_request",
					"30cc213cc6d97d6e256a850cdb778468cbd0c51d6cd62f406f84574c23d069ca",
				},
				signature: "ddfe3053ad8102a6a0843f660fe93bfa0adf85bc922c0a618a924e584e57b5a8",
				headers: []string{
					"PUT https://examplebucket.s3.amazonaws.com/photos/photo1.jpg HTTP/1.1",
					"Content-Length: 86",
					"Content-Encoding: aws-chunked",
					"X-Amz-Content-Sha256: STREAMING-AWS4-HMAC-SHA256-PAYLOAD",
					"X-Amz-Decoded-Content-Length: 0",
					"X-Amz-Date: 19840805T135000Z",
					"Authorization: AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/19840805/us-east-1/s3/aws4_request, SignedHeaders=content-encoding;content-length;host;x-amz-content-sha256;x-amz-date;x-amz-decoded-content-length, Signature=ddfe3053ad8102a6a0843f660fe93bfa0adf85bc922c0a618a924e584e57b5a8",
				},
				body: []string{
					"0;chunk-signature=3fbb6add718472d830071e1dede3d74f6f73ef41434d5c70ec5b0a1cb4ecbfb2\r\n",
					"\r\n",
				},
			},
		},
		{
			name: "STREAMING-AWS4-HMAC-SHA256-PAYLOAD/with payload",
			input: input{
				builder: func(req *fasthttp.Request) {
					req.Header.SetMethod("PUT")
					req.SetRequestURI(URI)
					req.SetBody(bytes.Repeat([]byte("a"), 64*1024+1024))
				},
				signer: &StreamedPayloadSigner{SignPayload: true},
			},
			expected: expected{
				canonicalRequest: []string{
					"PUT",
					"/photos/photo1.jpg",
					"",
					"content-encoding:aws-chunked",
					"content-length:66824",
					"host:examplebucket.s3.amazonaws.com",
					"x-amz-content-sha256:STREAMING-AWS4-HMAC-SHA256-PAYLOAD",
					"x-amz-date:19840805T135000Z",
					"x-amz-decoded-content-length:66560",
					"",
					"content-encoding;content-length;host;x-amz-content-sha256;x-amz-date;x-amz-decoded-content-length",
					"STREAMING-AWS4-HMAC-SHA256-PAYLOAD",
				},
				stringToSign: []string{
					"AWS4-HMAC-SHA256",
					"19840805T135000Z",
					"19840805/us-east-1/s3/aws4_request",
					"1e6b4e5faffa71a39b7a5032d2ab94599eaf16039d912de40d7f9531ddf26908",
				},
				signature: "7616b3ba4ff0a89d07292694e1ae0246235a64c48bfbefe2c4a69086cec02bb4",
				headers: []string{
					"PUT https://examplebucket.s3.amazonaws.com/photos/photo1.jpg HTTP/1.1",
					"Content-Length: 66824",
					"Content-Encoding: aws-chunked",
					"X-Amz-Content-Sha256: STREAMING-AWS4-HMAC-SHA256-PAYLOAD",
					"X-Amz-Decoded-Content-Length: 66560",
					"X-Amz-Date: 19840805T135000Z",
					"Authorization: AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/19840805/us-east-1/s3/aws4_request, SignedHeaders=content-encoding;content-length;host;x-amz-content-sha256;x-amz-date;x-amz-decoded-content-length, Signature=7616b3ba4ff0a89d07292694e1ae0246235a64c48bfbefe2c4a69086cec02bb4",
				},
				body: []string{
					"10000;chunk-signature=3d67bbd69c27b42bfb7ff17df70cc85150fe9d9d8fb264cf7acd6c9610e76afc\r\n",
					strings.Repeat("a", 64*1024) + "\r\n",
					"400;chunk-signature=b228a585c4757425775a64240fbba148e809a0b45641e238e1669fe24c772bbc\r\n",
					strings.Repeat("a", 1024) + "\r\n",
					"0;chunk-signature=4d8ffc8eb7d35ab8bca71a2a7f2c87a5ab5e0a2252a63628c3aecd616da56b5f\r\n",
					"\r\n",
				},
			},
		},
		{
			name: "STREAMING-AWS4-HMAC-SHA256-PAYLOAD-TRAILER/no payload/with trailer",
			input: input{
				builder: func(req *fasthttp.Request) {
					req.Header.SetMethod("PUT")
					req.SetRequestURI(URI)
					req.Header.Set(api.HeaderXAmzTrailer, api.HeaderXAmzChecksumCrc64nvme)
					req.Header.Set(api.HeaderXAmzChecksumCrc64nvme, "AAAAAAAAAAA=")
				},
				signer: &StreamedPayloadSigner{SignPayload: true},
			},
			expected: expected{
				canonicalRequest: []string{
					"PUT",
					"/photos/photo1.jpg",
					"",
					"content-encoding:aws-chunked",
					"content-length:215",
					"host:examplebucket.s3.amazonaws.com",
					"x-amz-content-sha256:STREAMING-AWS4-HMAC-SHA256-PAYLOAD-TRAILER",
					"x-amz-date:19840805T135000Z",
					"x-amz-decoded-content-length:0",
					"x-amz-trailer:x-amz-checksum-crc64nvme",
					"",
					"content-encoding;content-length;host;x-amz-content-sha256;x-amz-date;x-amz-decoded-content-length;x-amz-trailer",
					"STREAMING-AWS4-HMAC-SHA256-PAYLOAD-TRAILER",
				},
				stringToSign: []string{
					"AWS4-HMAC-SHA256",
					"19840805T135000Z",
					"19840805/us-east-1/s3/aws4_request",
					"77f2fdfdfbea9c064e3e26fc553dfcf499c2c90f4a5991caa8760f8d0ba40040",
				},
				signature: "7ed9db55e52a03fdedbcab8baad3590a784894e7058d3bbce6837cd9ea489761",
				headers: []string{
					"PUT https://examplebucket.s3.amazonaws.com/photos/photo1.jpg HTTP/1.1",
					"Content-Length: 215",
					"X-Amz-Trailer: x-amz-checksum-crc64nvme",
					"Content-Encoding: aws-chunked",
					"X-Amz-Content-Sha256: STREAMING-AWS4-HMAC-SHA256-PAYLOAD-TRAILER",
					"X-Amz-Decoded-Content-Length: 0",
					"X-Amz-Date: 19840805T135000Z",
					"Authorization: AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/19840805/us-east-1/s3/aws4_request, SignedHeaders=content-encoding;content-length;host;x-amz-content-sha256;x-amz-date;x-amz-decoded-content-length;x-amz-trailer, Signature=7ed9db55e52a03fdedbcab8baad3590a784894e7058d3bbce6837cd9ea489761",
				},
				body: []string{
					"0;chunk-signature=55118c8a4a8e532494a55ffb4947287876ce5c9a759cf4b3f559214341ade635\r\n",
					"x-amz-checksum-crc64nvme:AAAAAAAAAAA=\r\n",
					"x-amz-trailer-signature:e9269042d8a44972829ea8a325fa410a127f093dcb9f839bfc0a53e56b4c3809\r\n",
					"\r\n",
				},
			},
		},
		{
			name: "STREAMING-AWS4-HMAC-SHA256-PAYLOAD-TRAILER/no payload/empty trailer",
			input: input{
				builder: func(req *fasthttp.Request) {
					req.Header.SetMethod("PUT")
					req.SetRequestURI(URI)
				},
				signer: &StreamedPayloadSigner{SignPayload: true, ForceEmptyTrailer: true},
			},
			expected: expected{
				canonicalRequest: []string{
					"PUT",
					"/photos/photo1.jpg",
					"",
					"content-encoding:aws-chunked",
					"content-length:176",
					"host:examplebucket.s3.amazonaws.com",
					"x-amz-content-sha256:STREAMING-AWS4-HMAC-SHA256-PAYLOAD-TRAILER",
					"x-amz-date:19840805T135000Z",
					"x-amz-decoded-content-length:0",
					"",
					"content-encoding;content-length;host;x-amz-content-sha256;x-amz-date;x-amz-decoded-content-length",
					"STREAMING-AWS4-HMAC-SHA256-PAYLOAD-TRAILER",
				},
				stringToSign: []string{
					"AWS4-HMAC-SHA256",
					"19840805T135000Z",
					"19840805/us-east-1/s3/aws4_request",
					"b292793967a6eb8d04a8d7073b4a5d2273406f2b5ec7739ad1f9f5bb27ac5eb2",
				},
				signature: "e99cfc10645e62bd4fd9c2e988a876b740a63979fa251d6121940fa0cbe447a7",
				headers: []string{
					"PUT https://examplebucket.s3.amazonaws.com/photos/photo1.jpg HTTP/1.1",
					"Content-Length: 176",
					"Content-Encoding: aws-chunked",
					"X-Amz-Content-Sha256: STREAMING-AWS4-HMAC-SHA256-PAYLOAD-TRAILER",
					"X-Amz-Decoded-Content-Length: 0",
					"X-Amz-Date: 19840805T135000Z",
					"Authorization: AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/19840805/us-east-1/s3/aws4_request, SignedHeaders=content-encoding;content-length;host;x-amz-content-sha256;x-amz-date;x-amz-decoded-content-length, Signature=e99cfc10645e62bd4fd9c2e988a876b740a63979fa251d6121940fa0cbe447a7",
				},
				body: []string{
					"0;chunk-signature=da1798414de258fd5d494d5febfa27ea49e66e611cdb7dbc132a870b854b02e3\r\n",
					"x-amz-trailer-signature:9b9bd0c9114c4cf2e9d78293ce27f0ded13fc47f613f72225b612f8c35877a5b\r\n",
					"\r\n",
				},
			},
		},
		{
			name: "STREAMING-AWS4-HMAC-SHA256-PAYLOAD-TRAILER/with payload/with trailer",
			input: input{
				builder: func(req *fasthttp.Request) {
					req.Header.SetMethod("PUT")
					req.SetRequestURI(URI)
					req.SetBody(bytes.Repeat([]byte("a"), 64*1024+1024))
					req.Header.Set(api.HeaderXAmzTrailer, api.HeaderXAmzChecksumCrc32c)
					req.Header.Set(api.HeaderXAmzChecksumCrc32c, "sOO8/Q==")
				},
				signer: &StreamedPayloadSigner{SignPayload: true},
			},
			expected: expected{
				canonicalRequest: []string{
					"PUT",
					"/photos/photo1.jpg",
					"",
					"content-encoding:aws-chunked",
					"content-length:66946",
					"host:examplebucket.s3.amazonaws.com",
					"x-amz-content-sha256:STREAMING-AWS4-HMAC-SHA256-PAYLOAD-TRAILER",
					"x-amz-date:19840805T135000Z",
					"x-amz-decoded-content-length:66560",
					"x-amz-trailer:x-amz-checksum-crc32c",
					"",
					"content-encoding;content-length;host;x-amz-content-sha256;x-amz-date;x-amz-decoded-content-length;x-amz-trailer",
					"STREAMING-AWS4-HMAC-SHA256-PAYLOAD-TRAILER",
				},
				stringToSign: []string{
					"AWS4-HMAC-SHA256",
					"19840805T135000Z",
					"19840805/us-east-1/s3/aws4_request",
					"9c1af4363964ddc55e87be57e0c0bd04a6ad6169fd907110e5305b68a61c4b7b",
				},
				signature: "0f85426829ea47662c2fd47f8e27a38d52286d837acfc5e39c4f80e0e7d7540c",
				headers: []string{
					"PUT https://examplebucket.s3.amazonaws.com/photos/photo1.jpg HTTP/1.1",
					"Content-Length: 66946",
					"X-Amz-Trailer: x-amz-checksum-crc32c",
					"Content-Encoding: aws-chunked",
					"X-Amz-Content-Sha256: STREAMING-AWS4-HMAC-SHA256-PAYLOAD-TRAILER",
					"X-Amz-Decoded-Content-Length: 66560",
					"X-Amz-Date: 19840805T135000Z",
					"Authorization: AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/19840805/us-east-1/s3/aws4_request, SignedHeaders=content-encoding;content-length;host;x-amz-content-sha256;x-amz-date;x-amz-decoded-content-length;x-amz-trailer, Signature=0f85426829ea47662c2fd47f8e27a38d52286d837acfc5e39c4f80e0e7d7540c",
				},
				body: []string{
					"10000;chunk-signature=1fdb36871c0b95462d2a7e8ec2957eab97d686e2922dc940bc14640b5d4428f2\r\n",
					strings.Repeat("a", 64*1024) + "\r\n",
					"400;chunk-signature=9d627a83f34a6a5cff0f60f154e2a1e7c656cc861f1de5679b81a815b3859f7f\r\n",
					strings.Repeat("a", 1024) + "\r\n",
					"0;chunk-signature=47f2053bc6fd5e325c6e3f526ef7ddf62bc609b727a3d7a5f4d698375f4d2d70\r\n",
					"x-amz-checksum-crc32c:sOO8/Q==\r\n",
					"x-amz-trailer-signature:6dd0d219613b356ff545e69997fcabef2f491cb4f28d669b936c686c6cbb0490\r\n",
					"\r\n",
				},
			},
		},
		{
			name: "STREAMING-AWS4-HMAC-SHA256-PAYLOAD-TRAILER/with payload/empty trailer",
			input: input{
				builder: func(req *fasthttp.Request) {
					req.Header.SetMethod("PUT")
					req.SetRequestURI(URI)
					req.SetBody(bytes.Repeat([]byte("a"), 64*1024+1024))
				},
				signer: &StreamedPayloadSigner{SignPayload: true, ForceEmptyTrailer: true},
			},
			expected: expected{
				canonicalRequest: []string{
					"PUT",
					"/photos/photo1.jpg",
					"",
					"content-encoding:aws-chunked",
					"content-length:66914",
					"host:examplebucket.s3.amazonaws.com",
					"x-amz-content-sha256:STREAMING-AWS4-HMAC-SHA256-PAYLOAD-TRAILER",
					"x-amz-date:19840805T135000Z",
					"x-amz-decoded-content-length:66560",
					"",
					"content-encoding;content-length;host;x-amz-content-sha256;x-amz-date;x-amz-decoded-content-length",
					"STREAMING-AWS4-HMAC-SHA256-PAYLOAD-TRAILER",
				},
				stringToSign: []string{
					"AWS4-HMAC-SHA256",
					"19840805T135000Z",
					"19840805/us-east-1/s3/aws4_request",
					"ff8de5594bc36bfb6238ed96b78d31dcc8ef1d90aad92ee0b64a67025d7221b8",
				},
				signature: "7c172f5c0d74dea5eaaa30ea2bbd8ceb8982d45d764b70ee7b6ac1b5b0ff6999",
				headers: []string{
					"PUT https://examplebucket.s3.amazonaws.com/photos/photo1.jpg HTTP/1.1",
					"Content-Length: 66914",
					"Content-Encoding: aws-chunked",
					"X-Amz-Content-Sha256: STREAMING-AWS4-HMAC-SHA256-PAYLOAD-TRAILER",
					"X-Amz-Decoded-Content-Length: 66560",
					"X-Amz-Date: 19840805T135000Z",
					"Authorization: AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/19840805/us-east-1/s3/aws4_request, SignedHeaders=content-encoding;content-length;host;x-amz-content-sha256;x-amz-date;x-amz-decoded-content-length, Signature=7c172f5c0d74dea5eaaa30ea2bbd8ceb8982d45d764b70ee7b6ac1b5b0ff6999",
				},
				body: []string{
					"10000;chunk-signature=cee8f98b6d5097c9d053d11875eaa5073689e6536fbd76b03f445618e21a0532\r\n",
					strings.Repeat("a", 64*1024) + "\r\n",
					"400;chunk-signature=b0ae86e932a49b8f8bf39caacb8ffe14e4fd24504a5a3c2dff51550486c75462\r\n",
					strings.Repeat("a", 1024) + "\r\n",
					"0;chunk-signature=82aee475e413c0e62811982aef83095bea6aa73f22a11176bfaa95501c8d340c\r\n",
					"x-amz-trailer-signature:767264d4d943d876981186240c7be192b018d0257baddf233bf94b16200c94e0\r\n",
					"\r\n",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			signingTime := signing.SigningTimeOf(time.Date(1984, time.August, 5, 13, 50, 0, 0, time.UTC))

			var actualRequest fasthttp.Request
			actualRequest.Header.SetNoDefaultContentType(true)
			tc.input.builder(&actualRequest)

			canonicalRequest, stringToSign, signature, err := tc.input.signer.Sign(signing.SigningArgs{
				Request: &actualRequest,
				Credentials: signing.Credentials{
					AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
					SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
				},
				Region:      "us-east-1",
				SigningTime: signingTime,
			})
			require.NoError(t, err)

			assert.Equal(t, strings.Join(tc.expected.canonicalRequest, "\n"), canonicalRequest, "canonical request")
			assert.Equal(t, strings.Join(tc.expected.stringToSign, "\n"), stringToSign, "string to sign")
			assert.Equal(t, strings.Join(tc.expected.headers, "\r\n")+"\r\n\r\n", actualRequest.Header.String(), "headers")

			assert.Equal(t, strings.Join(tc.expected.body, ""), string(actualRequest.Body()), "body")
			assert.Equal(t, tc.expected.signature, signature, "signature")
		})
	}
}
