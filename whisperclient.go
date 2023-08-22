package whisperclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

const (
	audioTranscriptionURL string = "https://api.openai.com/v1/audio/transcriptions"
	FormatSrt             string = "srt"
	FormatText            string = "text"
	LanguageEnglish       string = "en"
	LanguagePortuguese    string = "pt"
)

// Client is a wrapper around the WhisperAI API
type Client struct {
	httpCli *http.Client
	apiKey  string
	model   string
}

// New returns a new Client
func New(httpCli *http.Client, apiKey, model string) *Client {
	return &Client{
		httpCli: httpCli,
		apiKey:  apiKey,
		model:   model,
	}
}

// TranscribeAudioInput is the input for the TranscribeAudio method
type TranscribeAudioInput struct {
	Name     string
	Language string
	Format   string
	Data     io.Reader
}

// TranscribeAudio transcribes the audio from the given input
func (c *Client) TranscribeAudio(ctx context.Context, in TranscribeAudioInput) ([]byte, error) {
	var body bytes.Buffer

	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", in.Name)
	if err != nil {
		return nil, fmt.Errorf("could not create form file: %w", err)
	}

	if _, err := io.Copy(part, in.Data); err != nil {
		return nil, fmt.Errorf("could not copy data to form file: %w", err)
	}

	if err := writer.WriteField("model", c.model); err != nil {
		return nil, fmt.Errorf("could not write model field: %w", err)
	}

	if err := writer.WriteField("language", in.Language); err != nil {
		return nil, fmt.Errorf("could not write language field: %w", err)
	}

	if err := writer.WriteField("response_format", in.Format); err != nil {
		return nil, fmt.Errorf("could not write response_format field: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("could not close writer: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, audioTranscriptionURL, &body)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	request.Header.Set("Authorization", "Bearer "+c.apiKey)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	response, err := c.httpCli.Do(request)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}
	defer response.Body.Close()

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}

	return b, nil
}
