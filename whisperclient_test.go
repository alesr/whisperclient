package whisperclient

import (
	"bytes"
	"context"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRoundTripper struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

func TestTranscribeAudio_SendsExpectedParameters(t *testing.T) {
	httpClient := http.Client{
		Transport: &mockRoundTripper{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, "Bearer test_api_key", req.Header.Get("Authorization"))

				contentType := req.Header.Get("Content-Type")
				assert.NotEmpty(t, contentType)

				_, params, err := mime.ParseMediaType(contentType)
				require.NoError(t, err, "Failed to parse content type")
				require.Contains(t, params, "boundary")

				reader := multipart.NewReader(req.Body, params["boundary"])

				parts := make(map[string]string)
				for {
					part, err := reader.NextPart()
					if err == io.EOF {
						break
					}
					require.NoError(t, err)

					partData, err := io.ReadAll(part)
					require.NoError(t, err)

					parts[part.FormName()] = string(partData)
				}

				assert.Equal(t, "test_model", parts["model"])
				assert.Equal(t, "en", parts["language"])
				assert.Equal(t, "text", parts["response_format"])

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer([]byte("mock response"))),
				}, nil
			},
		},
	}

	client := New(&httpClient, "test_api_key", "test_model")
	_, err := client.TranscribeAudio(context.TODO(), TranscribeAudioInput{
		Name:     "test_name",
		Language: LanguageEnglish,
		Format:   FormatText,
		Data:     bytes.NewBuffer([]byte("test audio data")),
	})

	require.NoError(t, err, "unexpected error")
}

func TestTranscribeAudio_ReturnsEquivalentData(t *testing.T) {
	expectedResponse := []byte("mock response")

	httpClient := http.Client{
		Transport: &mockRoundTripper{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(expectedResponse)),
				}, nil
			},
		},
	}

	client := New(&httpClient, "test_api_key", "test_model")

	response, err := client.TranscribeAudio(context.TODO(), TranscribeAudioInput{
		Name:     "test_name",
		Language: LanguageEnglish,
		Format:   FormatText,
		Data:     bytes.NewBuffer([]byte("test audio data")),
	})

	require.NoError(t, err, "unexpected error")
	assert.Equal(t, expectedResponse, response)
}
