package v4

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/lvjp/s3hobby/pkg/s3/signing"

	"github.com/stretchr/testify/require"
)

func TestStreamPayloadSigner(t *testing.T) {
	type testCaseInput struct {
		payload string
	}

	type testCaseExpected struct {
		stringToSign []string
		signature    string
	}

	chunks := []struct {
		input    testCaseInput
		expected testCaseExpected
	}{
		{
			input: testCaseInput{
				payload: strings.Repeat("a", 64*1024),
			},
			expected: testCaseExpected{
				stringToSign: []string{
					"AWS4-HMAC-SHA256-PAYLOAD",
					"19840805T135000Z",
					"19840805/eu-west-3/s3/aws4_request",
					"d81f82fc3505edab99d459891051a732e8730629a2e4a59689829ca17fe2e435",
					"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
					"bf718b6f653bebc184e1479f1935b8da974d701b893afcf49e701f3e2f9f9c5a",
				},
				signature: "4cf1e5fb361fad626520acbae58b18ed49deb0620c716090f4e39ffcde0f9cbe",
			},
		},
		{
			input: testCaseInput{
				payload: strings.Repeat("a", 1024),
			},
			expected: testCaseExpected{
				stringToSign: []string{
					"AWS4-HMAC-SHA256-PAYLOAD",
					"19840805T135000Z",
					"19840805/eu-west-3/s3/aws4_request",
					"4cf1e5fb361fad626520acbae58b18ed49deb0620c716090f4e39ffcde0f9cbe",
					"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
					"2edc986847e209b4016e141a6dc8716d3207350f416969382d431539bf292e4a",
				},
				signature: "3b77a0bcfc3f132e3d65aafc583ea951d3a5322c36fe67c14cb1d5e013520620",
			},
		},
		{
			expected: testCaseExpected{
				stringToSign: []string{
					"AWS4-HMAC-SHA256-PAYLOAD",
					"19840805T135000Z",
					"19840805/eu-west-3/s3/aws4_request",
					"3b77a0bcfc3f132e3d65aafc583ea951d3a5322c36fe67c14cb1d5e013520620",
					"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
					"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				},
				signature: "4ebbf68a5497487023968fdffa4c1b289aab3f6b50b8cc5b7dc3a845f3fd8cdb",
			},
		},
	}

	region := "eu-west-3"
	signingTime := signing.SigningTimeOf(time.Date(1984, time.August, 5, 13, 50, 0, 0, time.UTC))

	streamSigner := NewStreamPayloadSigner(
		NewSigningKey(
			signing.Credentials{
				AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
				SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			},
			region,
			signingTime,
		),
		"d81f82fc3505edab99d459891051a732e8730629a2e4a59689829ca17fe2e435",
		NewStringToSignBuilder(signingTime, region),
	)

	for i, chunk := range chunks {
		t.Run(fmt.Sprintf("chunk#%d", i), func(t *testing.T) {
			stringToSign, signature := streamSigner.ChunkSignature([]byte(chunk.input.payload))
			require.Equal(t, strings.Join(chunk.expected.stringToSign, "\n"), stringToSign, "stringToSign")
			require.Equal(t, chunk.expected.signature, signature, "signature")
		})

		if t.Failed() {
			t.Fatalf("cannot continue test, chunk %d failed", i)
			return
		}
	}

	t.Run("trailer", func(t *testing.T) {
		input := testCaseInput{
			payload: "x-amz-checksum-crc32c:sOO8/Q==\n",
		}
		expected := testCaseExpected{
			stringToSign: []string{
				"AWS4-HMAC-SHA256-TRAILER",
				"19840805T135000Z",
				"19840805/eu-west-3/s3/aws4_request",
				"4ebbf68a5497487023968fdffa4c1b289aab3f6b50b8cc5b7dc3a845f3fd8cdb",
				"1e376db7e1a34a8ef1c4bcee131a2d60a1cb62503747488624e10995f448d774",
			},
			signature: "e13314fde2b7451e34c30508968fca79fda6c44e7b63c86528e5a1e9b7615ce1",
		}

		stringToSign, signature := streamSigner.TrailerSignature([]byte(input.payload))
		require.Equal(t, strings.Join(expected.stringToSign, "\n"), stringToSign, "stringToSign")
		require.Equal(t, expected.signature, signature, "signature")
	})
}
