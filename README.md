# whisperclient
A simple Go client for interacting with the WhisperAI API.

## Overview

whisperclient is a Go package that provides a straightforward way to transcribe audio using the WhisperAI API.

## Installation

```bash
go get github.com/alesr/whisperclient
```

## Usage

```go
import "github.com/alesr/whisperclient"

// Initialize the client
client := whisperclient.New(httpClient, apiKey, model)

// Transcribe audio
input := whisperclient.TranscribeAudioInput{
    Name: "sample_audio",
    Language: whisperclient.LanguageEnglish,
    Format: whisperclient.FormatText,
    Data: audioDataReader,
}

response, err := client.TranscribeAudio(context.Background(), input)
```
